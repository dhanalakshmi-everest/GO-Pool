package Go_Pool

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	pool1 := NewPool(10, 2)
	pool2 := NewPool(10, -10)

	assert.Equal(t, pool1, pool2, "NewPool should return singleton struct")
	assert.Equal(t, 2, pool1.numCPUs, "No of CPUs used by pool should be 2")
}

func TestSetNumCPUs(t *testing.T) {
	pool := NewPool(10, 4)

	assert.Equal(t, 4, pool.numCPUs, "No of CPUs used by pool should be 4")

	pool.SetNumCPUs(2)
	assert.Equal(t, 2, pool.numCPUs, "No of CPUs used by pool should be 2")
}

func TestPool_Add(t *testing.T) {
	task1 := NewTask(
		func() (interface{}, error) {
			return "Task completed", nil
		}, time.Minute,
	)

	pool := NewPool(10, 4)

	pool.AddTask(task1)
	assert.Equal(t, "Task completed", <-task1.Result, "Task should be completed")

	task2 := NewTask(
		func() (interface{}, error) {
			return "Task completed", nil
		}, 0,
	)
	pool.AddTask(task2)
	assert.Equal(t, nil, <-task2.Result, "Task result should be nil")

	//multiple tasks
	var tasks = []*Task{
		NewTask(func() (interface{}, error) {
			time.Sleep(time.Second * 3)
			return "Task1 completed", nil
		}, time.Minute),

		NewTask(func() (interface{}, error) {
			time.Sleep(time.Second * 2)
			return "Task2 completed", nil
		}, time.Minute),

		NewTask(func() (interface{}, error) {
			time.Sleep(time.Second * 1)
			return "Task3 completed", nil
		}, time.Minute),
	}

	for _, task := range tasks {
		pool.AddTask(task)
	}

	assert.Equal(t, "Task3 completed", <-tasks[2].Result, "Task3 should be completed")
	assert.Equal(t, "Task2 completed", <-tasks[1].Result, "Task2 should be completed")
	assert.Equal(t, "Task1 completed", <-tasks[0].Result, "Task1 should be completed")

}
