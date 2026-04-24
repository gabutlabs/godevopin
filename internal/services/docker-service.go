package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// DockerService handles Docker operations
type DockerService struct {
	client *client.Client
}

// ContainerInfo represents container information
type ContainerInfo struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Image          string            `json:"image"`
	Status         string            `json:"status"`
	State          string            `json:"state"`
	Labels         map[string]string `json:"labels"`
	Summary        container.Summary `json:"summary"`
	ComposeProject string            `json:"compose_project,omitempty"`
}

type ContainerDetailInfo struct {
	ContainerInfo   `json:"container_info"`
	Ports           []container.Port                  `json:"ports"`
	NetworkSettings *container.NetworkSettingsSummary `json:"network_settings_summary"`
	Mounts          []container.MountPoint            `json:"mounts"`
}

// ImageInfo represents image information
type ImageInfo struct {
	ID      string   `json:"id"`
	Tags    []string `json:"tags"`
	SizeMB  float64  `json:"size_mb"`
	Created int64    `json:"created"`
}

// NetworkInfo represents network information
type NetworkInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Driver string `json:"driver"`
	Scope  string `json:"scope"`
}

// VolumeInfo represents volume information
type VolumeInfo struct {
	Name       string `json:"name"`
	Driver     string `json:"driver"`
	Mountpoint string `json:"mountpoint"`
}

// NewDockerService creates a new Docker service instance
func NewDockerService() *DockerService {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("failed to create docker client: %v", err)
		return nil
	}

	return &DockerService{
		client: cli,
	}
}

// Close closes the Docker client connection
func (s *DockerService) Close() error {
	return s.client.Close()
}

// GetContainers retrieves all containers
func (s *DockerService) GetContainers(ctx context.Context, all bool) ([]ContainerInfo, error) {
	containers, err := s.client.ContainerList(ctx, container.ListOptions{All: all})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	result := make([]ContainerInfo, 0, len(containers))
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		result = append(result, ContainerInfo{
			ID:             c.ID[:12],
			Name:           name,
			Image:          c.Image,
			Status:         c.Status,
			State:          c.State,
			Labels:         c.Labels,
			Summary:        c,
			ComposeProject: c.Labels["com.docker.compose.project"],
		})
	}

	return result, nil
}

// GetContainers retrieves all containers
func (s *DockerService) GetContainerInfo(ctx context.Context, id string) (*ContainerDetailInfo, error) {
	filter := filters.NewArgs()
	filter.Add("id", id)
	containers, err := s.client.ContainerList(ctx, container.ListOptions{Filters: filter, All: true})

	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("container with ID %q not found", id)
	}

	c := containers[0]

	name := ""
	if len(c.Names) > 0 {
		name = strings.TrimPrefix(c.Names[0], "/")
	}

	// Ambil 12 karakter pertama dari ID (aman)
	shortID := c.ID
	if len(shortID) > 12 {
		shortID = shortID[:12]
	}

	result := ContainerDetailInfo{
		ContainerInfo: ContainerInfo{
			ID:     shortID,
			Name:   name,
			Image:  c.Image,
			Status: c.Status,
			State:  c.State,
			Labels: c.Labels,
		},
		NetworkSettings: c.NetworkSettings,
		Mounts:          c.Mounts,
		Ports:           c.Ports,
	}
	return &result, nil
}

func (s *DockerService) GetContainerLogs(ctx context.Context, id string, tail int) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       fmt.Sprintf("%d", tail),
		Timestamps: true, // opsional: tambahin timestamp kalau mau
	}

	logsReader, err := s.client.ContainerLogs(ctx, id, options)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logsReader.Close()

	// Gunakan docker's own demux function supaya bersih
	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, logsReader)
	if err != nil {
		return "", fmt.Errorf("failed to demux logs: %w", err)
	}

	// Gabungkan stdout dan stderr (biasanya stdout sudah cukup, tapi tergantung app)
	// Kalau mau bedain, bisa return terpisah atau tambah prefix [STDOUT]/[STDERR]
	var result strings.Builder
	if stdoutBuf.Len() > 0 {
		result.WriteString(stdoutBuf.String())
	}
	if stderrBuf.Len() > 0 {
		// Opsional: tambah prefix biar kelihatan ini error
		lines := strings.Split(stderrBuf.String(), "\n")
		for _, line := range lines {
			if line != "" {
				result.WriteString("[STDERR] ")
				result.WriteString(line)
				result.WriteString("\n")
			}
		}
	}

	return result.String(), nil
}

// StreamContainerLogs streams container logs in real-time
func (s *DockerService) StreamContainerLogs(ctx context.Context, id string, lines chan<- string) error {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false, // kalau mau timestamp, bisa false kalau tidak
		Tail:       "500", // optional: kalau mau semua historical logs dulu
	}

	logsReader, err := s.client.ContainerLogs(ctx, id, options)
	if err != nil {
		return fmt.Errorf("failed to get container logs: %w", err)
	}

	scanner := bufio.NewScanner(logsReader)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) < 8 {
			return 0, nil, nil // need more data
		}

		// Read header
		streamType := data[0] // 1 = stdout, 2 = stderr
		size := binary.BigEndian.Uint32(data[4:8])

		if len(data) < 8+int(size) {
			return 0, nil, nil // need more data
		}

		payload := data[8 : 8+size]

		// Prefix untuk distinguish stderr
		var prefixed []byte
		if streamType == 2 { // stderr
			prefixed = append([]byte("[STDERR] "), payload...)
		} else {
			prefixed = payload
		}

		return 8 + int(size), prefixed, nil
	})

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			logsReader.Close()
			return ctx.Err()
		default:
			line := scanner.Text()
			if strings.TrimSpace(line) != "" {
				select {
				case lines <- line:
				case <-ctx.Done():
					logsReader.Close()
					return ctx.Err()
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logsReader.Close()
		return fmt.Errorf("error reading logs: %w", err)
	}

	logsReader.Close()
	return nil
}

func (s *DockerService) ExecInteractive(ctx context.Context, containerID string, cmd []string) (types.HijackedResponse, string, error) {
	// Config untuk interactive shell
	execConfig := container.ExecOptions{
		AttachStdin:  true, // Penting untuk input
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true, // TTY mode - ini yang bikin interactive!
		Cmd:          cmd,
	}

	// Create exec instance
	execIDResp, err := s.client.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return types.HijackedResponse{}, "", fmt.Errorf("failed to create exec instance: %w", err)
	}

	// Attach dengan Tty=true
	resp, err := s.client.ContainerExecAttach(ctx, execIDResp.ID, container.ExecAttachOptions{
		Tty: true, // Penting!
	})
	if err != nil {
		return types.HijackedResponse{}, "", fmt.Errorf("failed to attach to exec instance: %w", err)
	}
	err = s.client.ContainerExecStart(ctx, execIDResp.ID, container.ExecStartOptions{
		Tty: true,
	})
	if err != nil {
		resp.Close()
		return types.HijackedResponse{}, "", fmt.Errorf("failed to start exec instance: %w", err)
	}

	return resp, execIDResp.ID, nil
}

// ResizeExec - untuk resize terminal
func (s *DockerService) ResizeExec(ctx context.Context, execID string, height, width uint) error {
	return s.client.ContainerExecResize(ctx, execID, container.ResizeOptions{
		Height: height,
		Width:  width,
	})
}

// GetImages retrieves all images
func (s *DockerService) GetImages(ctx context.Context) ([]ImageInfo, error) {
	images, err := s.client.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	result := make([]ImageInfo, 0, len(images))
	for _, img := range images {
		tags := img.RepoTags
		if tags == nil {
			tags = []string{"<none>"}
		}

		result = append(result, ImageInfo{
			ID:      img.ID[7:19], // sha256:xxxx -> xxxx
			Tags:    tags,
			SizeMB:  math.Round(float64(img.Size)/1024/1024*100) / 100,
			Created: img.Created,
		})
	}

	return result, nil
}

// GetNetworks retrieves all networks
func (s *DockerService) GetNetworks(ctx context.Context) ([]NetworkInfo, error) {
	networks, err := s.client.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	result := make([]NetworkInfo, 0, len(networks))
	for _, nw := range networks {
		result = append(result, NetworkInfo{
			ID:     nw.ID[:12],
			Name:   nw.Name,
			Driver: nw.Driver,
			Scope:  nw.Scope,
		})
	}

	return result, nil
}

// GetVolumes retrieves all volumes
func (s *DockerService) GetVolumes(ctx context.Context) ([]VolumeInfo, error) {
	volumes, err := s.client.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	result := make([]VolumeInfo, 0, len(volumes.Volumes))
	for _, vol := range volumes.Volumes {
		result = append(result, VolumeInfo{
			Name:       vol.Name,
			Driver:     vol.Driver,
			Mountpoint: vol.Mountpoint,
		})
	}

	return result, nil
}

// GetAllResources retrieves all Docker resources at once
func (s *DockerService) GetAllResources(ctx context.Context) (map[string]interface{}, error) {
	containers, err := s.GetContainers(ctx, true)
	if err != nil {
		return nil, err
	}

	images, err := s.GetImages(ctx)
	if err != nil {
		return nil, err
	}

	networks, err := s.GetNetworks(ctx)
	if err != nil {
		return nil, err
	}

	volumes, err := s.GetVolumes(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"containers": containers,
		"images":     images,
		"networks":   networks,
		"volumes":    volumes,
	}, nil
}
