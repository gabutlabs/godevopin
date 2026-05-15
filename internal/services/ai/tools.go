package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gabutlabs/devopin/internal/model"
	service "github.com/gabutlabs/devopin/internal/services"
	"github.com/tmc/langchaingo/llms"
)

type AIToolProvider struct {
	systemService  service.SystemMetricService
	workerService  service.WorkerServiceService
	dockerService  *service.DockerService
	logService     service.LogHistoryService
	projectService service.ProjectService
}

func NewAIToolProvider(
	sys service.SystemMetricService,
	wrk service.WorkerServiceService,
	doc *service.DockerService,
	log service.LogHistoryService,
	prj service.ProjectService,
) *AIToolProvider {
	return &AIToolProvider{
		systemService:  sys,
		workerService:  wrk,
		dockerService:  doc,
		logService:     log,
		projectService: prj,
	}
}

// GetTools returns the list of tools for the AI agent
func (p *AIToolProvider) GetTools() []llms.Tool {
	return []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "get_system_metrics",
				Description: "Get current system metrics including CPU, Memory, and Disk usage.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"filter": map[string]any{
							"type":        "string",
							"description": "Filter range (e.g., '1h', '6h', '1d'). Default is '1h'.",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "manage_worker",
				Description: "Manage background worker services (list, start, stop, restart).",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"action": map[string]any{
							"type": "string",
							"enum": []string{"list", "start", "stop", "restart"},
						},
						"worker_name": map[string]any{
							"type":        "string",
							"description": "The name of the worker service (required for start/stop/restart).",
						},
					},
					"required": []string{"action"},
				},
			},
		},
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "manage_docker",
				Description: "Manage Docker containers (list, start, stop, restart).",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"action": map[string]any{
							"type": "string",
							"enum": []string{"list", "start", "stop", "restart"},
						},
						"container_id": map[string]any{
							"type":        "string",
							"description": "The ID or Name of the Docker container (required for start/stop/restart).",
						},
					},
					"required": []string{"action"},
				},
			},
		},
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "analyze_project_logs",
				Description: "Retrieve and analyze logs for a specific project.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"project_name": map[string]any{
							"type":        "string",
							"description": "The name of the project to analyze logs for.",
						},
						"limit": map[string]any{
							"type":        "integer",
							"description": "Number of log lines to retrieve (max 500).",
						},
					},
					"required": []string{"project_name"},
				},
			},
		},
	}
}

// ExecuteTool executes a tool call from the AI agent
func (p *AIToolProvider) ExecuteTool(ctx context.Context, toolCall llms.ToolCall) (string, error) {
	switch toolCall.FunctionCall.Name {
	case "get_system_metrics":
		var args struct {
			Filter string `json:"filter"`
		}
		if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
			return "", err
		}
		if args.Filter == "" {
			args.Filter = "1h"
		}
		metrics, err := p.systemService.GetSystemMetrics()
		if err != nil {
			return "", err
		}
		data, _ := json.Marshal(metrics)
		return string(data), nil

	case "manage_worker":
		var args struct {
			Action     string `json:"action"`
			WorkerName string `json:"worker_name"`
		}
		if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
			return "", err
		}

		if args.Action == "list" {
			workers, err := p.workerService.ListWorkerServices(nil)
			if err != nil {
				return "", err
			}
			data, _ := json.Marshal(workers)
			return string(data), nil
		}

		worker, err := p.workerService.GetWorkerServiceByName(args.WorkerName)
		if err != nil {
			return fmt.Sprintf("Worker '%s' not found", args.WorkerName), nil
		}

		var status model.CurrentStatus
		switch args.Action {
		case "start":
			status = model.StatusStarting
		case "stop":
			status = model.StatusStopped
		case "restart":
			status = model.StatusRestart
		default:
			return "Invalid action", nil
		}

		err = p.workerService.UpdateWorkerServiceStatus(worker.ID, status, worker.HealthStatus)
		if err != nil {
			return fmt.Sprintf("Failed to %s worker: %v", args.Action, err), nil
		}
		return fmt.Sprintf("Successfully triggered %s for worker %s", args.Action, args.WorkerName), nil

	case "manage_docker":
		var args struct {
			Action      string `json:"action"`
			ContainerID string `json:"container_id"`
		}
		if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
			return "", err
		}

		if args.Action == "list" {
			containers, err := p.dockerService.GetContainers(ctx, true)
			if err != nil {
				return "", err
			}
			data, _ := json.Marshal(containers)
			return string(data), nil
		}

		var err error
		switch args.Action {
		case "start":
			err = p.dockerService.StartContainer(ctx, args.ContainerID)
		case "stop":
			err = p.dockerService.StopContainer(ctx, args.ContainerID)
		case "restart":
			err = p.dockerService.RestartContainer(ctx, args.ContainerID)
		default:
			return "Invalid action", nil
		}

		if err != nil {
			return fmt.Sprintf("Failed to %s container: %v", args.Action, err), nil
		}
		return fmt.Sprintf("Successfully triggered %s for container %s", args.Action, args.ContainerID), nil

	case "analyze_project_logs":
		var args struct {
			ProjectName string `json:"project_name"`
			Limit       int    `json:"limit"`
		}
		if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
			return "", err
		}

		if args.Limit <= 0 || args.Limit > 500 {
			args.Limit = 100
		}

		projects, err := p.projectService.ListProjects()
		if err != nil {
			return "", err
		}

		var targetProject *model.Project
		for _, prj := range projects {
			if strings.EqualFold(prj.Name, args.ProjectName) {
				targetProject = &prj
				break
			}
		}

		if targetProject == nil {
			return fmt.Sprintf("Project '%s' not found", args.ProjectName), nil
		}

		logs, _, err := p.logService.GetLogsByProjectID(targetProject.ID, args.Limit, 0)
		if err != nil {
			return "", err
		}

		var logContents []string
		for _, l := range logs {
			logContents = append(logContents, fmt.Sprintf("[%s] %s: %s", l.Timestamp, l.Level, l.Message))
		}

		return strings.Join(logContents, "\n"), nil

	default:
		return "", fmt.Errorf("unknown tool: %s", toolCall.FunctionCall.Name)
	}
}
