package worker

import (
	"fmt"
	"cube/task"
	"github.com/google/uuid"
	"github.com/golang-collections/collections/queue"
)

type Worker struct {
	ID uuid.UUID
	Name string
	Db map[uuid.UUID]*task.Task
	Queue queue.Queue
	TaskCount int
}

func (w *Worker) RunTask(){
	fmt.Printf("I will run the task\n")
}

func (w *Worker) StartTask(){
	fmt.Printf("I will start the task\n")
}

func (w *Worker) StopTask(){
	fmt.Printf("I will Stop the task\n")
}

func (w *Worker) CollectStats(){
	fmt.Printf("I will collect status from the task\n")
}
