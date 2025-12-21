package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// DockerService handles Docker operations
type DockerService struct {
	client *client.Client
}

// ContainerInfo represents container information
type ContainerInfo struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	Status  string            `json:"status"`
	State   string            `json:"state"`
	Labels  map[string]string `json:"labels"`
	Summary container.Summary `json:"summary"`
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
			ID:      c.ID[:12],
			Name:    name,
			Image:   c.Image,
			Status:  c.Status,
			State:   c.State,
			Labels:  c.Labels,
			Summary: c,
		})
	}

	return result, nil
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
