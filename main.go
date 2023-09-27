package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
)

type ContainerInfo struct {
	ID  string
	PID uint32
}

func getRunningContainers() ([]ContainerInfo, error) {
	ctx := namespaces.WithNamespace(context.Background(), "moby")

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		fmt.Printf("Error connecting to Containerd: %v\n", err)
		return nil, err
	}
	defer client.Close()
	containers, err := client.Containers(ctx)
	if err != nil {
		fmt.Printf("Error listing containers: %v\n", err)
		return nil, err
	}

	var containerInfos []ContainerInfo

	for _, container := range containers {
		info, err := container.Info(ctx)
		if err != nil {
			fmt.Printf("Error getting container info: %v\n", err)
			continue
		}

		task, err := container.Task(ctx, cio.NewAttach())
		containerInfos = append(containerInfos, ContainerInfo{
			ID:  info.ID,
			PID: task.Pid(),
		})

		fmt.Printf("Container ID: %s. PID: %d\n", info.ID, task.Pid())
	}
	return containerInfos, nil
}

func createSymLinks(cn ContainerInfo) error {
	pidStr := strconv.Itoa(int(cn.PID))
	targetPath := "/proc/" + pidStr + "/ns/net"
	symlinkPath := "/var/run/netns/" + cn.ID

	err := os.Symlink(targetPath, symlinkPath)
	if err != nil {
		return nil
	}

	println("Symlink created:", symlinkPath)
	return nil
}

func main() {
	cns, err := getRunningContainers()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range cns {
		err := createSymLinks(c)
		if err != nil {
			fmt.Printf("Error processing container: %s, %s \n", c.ID, err.Error())
		}
	}
}
