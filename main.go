package main

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
)

type ContainerInfo struct {
	ID  string
	PID uint32
}

func getRunningContainers() {
	ctx := namespaces.WithNamespace(context.Background(), "moby")

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		fmt.Printf("Error connecting to Containerd: %v\n", err)
		return
	}
	defer client.Close()
	containers, err := client.Containers(ctx)
	if err != nil {
		fmt.Printf("Error listing containers: %v\n", err)
		return
	}

	var containerInfo []ContainerInfo

	for _, container := range containers {
		info, err := container.Info(ctx)
		if err != nil {
			fmt.Printf("Error getting container info: %v\n", err)
			continue
		}

		task, err := container.Task(ctx, cio.NewAttach())
		containerInfo = append(containerInfo, ContainerInfo{
			ID:  info.ID,
			PID: task.Pid(),
		})

		fmt.Printf("Container ID: %s. PID: %d\n", info.ID, task.Pid())
	}
}

func main() {
	getRunningContainers()
}
