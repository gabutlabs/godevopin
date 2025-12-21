package worker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type DockerManagement struct{}

func NewDockerManagement() *DockerManagement {
	return &DockerManagement{}
}

func (dm *DockerManagement) RunDocker() {
	ctx := context.Background()

	// Create docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// get all container like "docker ps -a"
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Containers ===")
	for _, c := range containers {
		fmt.Printf("ID: %s | Name: %s | Status: %s\n",
			c.ID[:12], c.Names[0], c.Status)
	}

	// get all images
	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("\n=== Images ===")
	for _, img := range images {
		if len(img.RepoTags) > 0 {
			fmt.Printf("Tag: %s | Size: %.2f MB\n",
				img.RepoTags[0], float64(img.Size)/1024/1024)
		}
	}

	network, err := cli.NetworkList(ctx, network.ListOptions{})
	fmt.Println("\n=== network ===")
	for _, nw := range network {
		fmt.Printf("network name: %s", nw.Name)
	}
}
