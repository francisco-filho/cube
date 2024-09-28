package main

import (
	"fmt"
	"time"
	"cube/task"
	"cube/manager"
	"cube/node"
	"cube/worker"

	"github.com/google/uuid"
	"github.com/golang-collections/collections/queue"
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

}
