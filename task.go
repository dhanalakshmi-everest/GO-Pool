package Go_Pool

import (
	"errors"
	"time"
)

var TaskMaxWaitingExceeded = errors.New("Task maximum waiting time exceeded")

type Functionality func() (interface{}, error)

type Task struct {
	functionality  Functionality
	arrivalTime    time.Time
	maxWaitingTime time.Duration
	Result         chan interface{}
	Error          error
}

func NewTask(fn Functionality, maxWaitingTime time.Duration) *Task {
	return &Task{
		functionality:  fn,
		arrivalTime:    time.Now(),
		maxWaitingTime: maxWaitingTime,
		Result:         make(chan interface{}, 1),
	}
}

func (task *Task) Run() {

	timeElapsed := time.Since(task.arrivalTime)

	if timeElapsed >= task.maxWaitingTime {
		task.Error = TaskMaxWaitingExceeded
		task.Result <- nil
		return
	}

	result, error := task.functionality()

	if error != nil {
		task.Error = error
		task.Result <- nil
	} else {
		task.Result <- result
	}

}
