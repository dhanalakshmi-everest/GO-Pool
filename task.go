package Go_Pool

import "sync"

type Task struct {
	Err error
	fn  func() error
}

func NewTask(f func() error) *Task {
	return &Task{fn: f}
}

func (task *Task) Run(wg *sync.WaitGroup) {
	task.Err = task.fn()
	wg.Done()
}
