package manager

import (
	"fmt"
	"cube/task"
	"github.com/google/uuid"
	"github.com/golang-collections/collections/queue"
)

type Manager struct {
	Pending queue.Queue
	TaskDb map[string][]*task.Task
	EventDb map[string][]*task.TaskEvent
	Workers []string
	WorkersTaskMap map[string][]uuid.UUID
	TasksWorkerMap map[uuid.UUID]string
}

func (m *Manager) SelectWorker(){
	fmt.Printf("SelectWorker\n")
}

func (m *Manager) CollectStats(){
	fmt.Printf("Managing collect stats\n")
}

func (m *Manager) UpdateTasks(){
	fmt.Printf("Update tasks\n")
}

func (m *Manager) SendWork(){
	fmt.Printf("Sending work...\n")
}
