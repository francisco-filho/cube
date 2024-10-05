package main

import (
	"os"
	"log"
	"fmt"
	"time"
	"cube/task"
	"cube/manager"
	"cube/node"
	"cube/worker"

	"github.com/google/uuid"
	"github.com/golang-collections/collections/queue"
	"github.com/docker/docker/client"
)

func main(){

	t := task.Task{
		ID: uuid.New(),
		Name: "task-1",
		State: task.Pending,
		Image: "ubuntu",
		Memory: 1024,
		Disk: 4096,
	}

	te := task.TaskEvent{
		ID: uuid.New(),
		State: task.Pending,
		Timestamp: time.Now(),
		Task: t,
	}

	fmt.Printf("task: %v\n", t)
	fmt.Printf("taskEvent: %v\n", te)

	w := worker.Worker{
		ID: uuid.New(),
		Name: "worker-1",
		Queue: *queue.New(),
		Db: make(map[uuid.UUID]*task.Task),
	}

	fmt.Printf("worker: %v\n", w)
	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		Pending: *queue.New(),
		TaskDb: make(map[string][]*task.Task),
		EventDb: make(map[string][]*task.TaskEvent),
		Workers: []string{w.Name},
	}

	fmt.Printf("Manager: %v\n", m)
	m.SelectWorker()
	m.UpdateTasks()
	m.SendWork()

    n := node.Node{
        Name:   "Node-1",
        Ip:     "192.168.1.1",
        Cores:  4,
        Memory: 1024,
        Disk:   25,
        Role:   "worker",
    }
 
    fmt.Printf("node: %v\n", n)

	fmt.Println("Create test container")
	dockerTask, createResult := createContainer()	
	if createResult.Error != nil {
		log.Printf("%v\n", createResult.Error)
		os.Exit(1)
	}

	time.Sleep(10 * time.Second)

	stopResult := stopContainer(dockerTask, createResult.ContainerId)
	if stopResult.Error != nil {
		log.Printf("%v\n", stopResult.Error)
		os.Exit(1)
	}

}


func createContainer()(*task.Docker, *task.DockerResult){
	c := task.Config{
		Name: "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=cube",
		},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv)

	docker := task.Docker{Client: dc, Config: c}	

	r := docker.Run()

	if r.Error != nil {
		log.Printf("Error Running Docker %v\n", r.Error)
		return nil, &r
	}

	log.Printf("Container %s is running with config %v\n", r.ContainerId, c)
	return &docker, &r
}

func stopContainer(d *task.Docker, id string) *task.DockerResult{
	res := d.Stop(id)
	if res.Error != nil {
		log.Printf("Error stoping the container %s: %v\n", id, res.Error)
		return &res
	}

	log.Printf("Container %s stopped\n", id)

	return &res
}


