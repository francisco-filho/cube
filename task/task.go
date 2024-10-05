package task

import (
	"log"
	"os"
	"math"
	"context"
	"io"
	"time"
	"github.com/google/uuid"
	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/container"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID uuid.UUID
	Name string
	State State
	Image string
	Memory int
	Disk int
	ExposedPorts nat.PortSet
	PortBidings map[string]string
	RestartPolicy string
	StartTime time.Time
	FinishTime time.Time
}

type TaskEvent struct {
	ID uuid.UUID
	State State
	Task Task
	Timestamp time.Time
}

type Runtime struct {
	ContainerId string
}

type Config struct {
	Name string
	AttachStdin	bool
	AttachStdout bool
	AttachStderr bool
	ExposedPorts nat.PortSet
	Cmd []string
	Image string
	Cpu float64
	Memory int64
	Disk int64
	Env []string
	RestartPolicy string
	Runtime Runtime
}

type Docker struct {
	Client *client.Client
	Config Config
}

type DockerResult struct {
	Error error
	Action string
	ContainerId string
	Result string
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()

	// Pull the image from dockerhub
	reader, err := d.Client.ImagePull(ctx, d.Config.Image, image.PullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}

	io.Copy(os.Stdout, reader)

	// Configuration
	rp := container.RestartPolicy{
		Name: container.RestartPolicyMode(d.Config.RestartPolicy),
	}

	r := container.Resources {
		Memory: d.Config.Memory,
		NanoCPUs: int64(d.Config.Cpu * math.Pow(10, 9)),
	}

	cc := container.Config{
		Image: d.Config.Image,
		Tty: false,
		Env: d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}

	hc := container.HostConfig{
		Resources: r,
		RestartPolicy: rp,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error CreateContainer %s: %v\n", d.Config.Image, err)
		return DockerResult{Error:err}
	}

	err = d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Printf("Error ContainerStart %s: %v\n", resp.ID, err)
		return DockerResult{Error:err, ContainerId: resp.ID}
	}

	d.Config.Runtime.ContainerId = resp.ID

out, err := d.Client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Printf("Error on logs: %s: %v", resp.ID, err)
		return DockerResult{Error:err, ContainerId: resp.ID}
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{ContainerId: resp.ID, Action: "Start", Error: nil}
}



