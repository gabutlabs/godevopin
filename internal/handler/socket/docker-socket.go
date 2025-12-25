package socket

import (
	"context"
	"log"

	service "github.com/gabutlabs/devopin/internal/services"
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
			log.Println(line)
			c.WriteMessage(websocket.TextMessage, []byte(line))

		case <-ctx.Done():
			log.Printf("Stream cancelled (client disconnect)")
			return
		}
	}
}

func (h *DockerSocketHandler) SetupDockerSocketRoutes(wsGroup fiber.Router) {
	wsGroup.Get("/docker/container/logs/:id", websocket.New(h.HandleContainerLogs))
}
