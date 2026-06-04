package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	service "github.com/gabutlabs/godevopin/internal/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type DockerSocketHandler struct {
	s *service.DockerService
}

func NewDockerSocketHandler(service *service.DockerService) *DockerSocketHandler {
	return &DockerSocketHandler{
		s: service,
	}
}

func (h *DockerSocketHandler) HandleContainerLogs(c *websocket.Conn) {
	containerID := c.Params("id")
	// UUID kalau butuh: bisa pakai c.Locals atau generate sendiri
	log.Printf("Client connected to stream logs for container: %s", containerID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logLines := make(chan string, 100)

	go func() {
		err := h.s.StreamContainerLogs(ctx, containerID, logLines)
		log.Println(err)
		if err != nil {
			log.Printf("Error streaming logs for container %s: %v", containerID, err)
		}
		close(logLines)
	}()

	for {
		select {
		case line, ok := <-logLines:
			if !ok {
				// End of logs
				log.Printf("End of log stream for container: %s", containerID)
				c.WriteMessage(websocket.TextMessage, []byte("--- End of logs ---"))
				return
			}
			c.WriteMessage(websocket.TextMessage, []byte(line))

		case <-ctx.Done():
			log.Printf("Stream cancelled (client disconnect)")
			return
		}
	}
}

func (h *DockerSocketHandler) HandlerContainerExec(c *websocket.Conn) {
	type ResizeMessage struct {
		Type string `json:"type"`
		Rows uint   `json:"rows"`
		Cols uint   `json:"cols"`
	}
	containerID := c.Params("id")
	commandQuery := c.Query("command", "/bin/bash")

	log.Printf("Client connected to interactive exec for container: %s", containerID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var command []string
	if commandQuery != "" {
		command = []string{commandQuery}
	} else {
		command = []string{"/bin/bash"}
	}

	// Start interactive exec
	resp, execID, err := h.s.ExecInteractive(ctx, containerID, command)
	if err != nil {
		log.Printf("Failed to start exec: %v", err)
		c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v\r\n", err)))
		return
	}
	defer resp.Close()

	// Kirim newline awal untuk trigger prompt
	resp.Conn.Write([]byte("\n"))

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1: Docker output -> WebSocket
	go func() {
		defer wg.Done()
		defer log.Println("Docker output reader stopped")

		buf := make([]byte, 8192)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := resp.Reader.Read(buf)
				if err != nil {
					if err != io.EOF {
						log.Printf("Error reading from docker: %v", err)
					}
					return
				}

				if n > 0 {
					if err := c.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
						log.Printf("Error writing to websocket: %v", err)
						return
					}
				}
			}
		}
	}()

	// Goroutine 2: WebSocket input -> Docker stdin
	go func() {
		defer wg.Done()
		defer log.Println("WebSocket input reader stopped")
		defer cancel()

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Println("WebSocket closed by client")
				} else {
					log.Printf("Error reading from websocket: %v", err)
				}
				return
			}

			// Handle resize command
			if messageType == websocket.TextMessage {
				var resizeMsg ResizeMessage
				if err := json.Unmarshal(message, &resizeMsg); err == nil {
					if resizeMsg.Type == "resize" {
						log.Printf("Resizing terminal to %dx%d", resizeMsg.Cols, resizeMsg.Rows)
						if err := h.s.ResizeExec(ctx, execID, resizeMsg.Rows, resizeMsg.Cols); err != nil {
							log.Printf("Failed to resize: %v", err)
						}
						continue
					}
				}
			}

			// Send raw input to Docker stdin
			_, err = resp.Conn.Write(message)
			if err != nil {
				log.Printf("Error writing to docker stdin: %v", err)
				return
			}
		}
	}()

	wg.Wait()
	log.Printf("Interactive exec session ended for container: %s", containerID)
}

func (h *DockerSocketHandler) SetupDockerSocketRoutes(wsGroup fiber.Router) {
	wsGroup.Get("/docker/container/logs/:id", websocket.New(h.HandleContainerLogs))
	wsGroup.Get("/docker/container/exec/:id", websocket.New(h.HandlerContainerExec))
}
