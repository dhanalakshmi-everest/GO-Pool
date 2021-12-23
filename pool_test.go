package Go_Pool

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var results = make(chan string, 3)
var tasks = []*Task{
	{
		fn: func() error {
			time.Sleep(time.Second * 3)
			results <- "Task1 completed"
			return nil
		}},
	{
		fn: func() error {
			time.Sleep(time.Second * 2)
			results <- "Task2 completed"
			return nil
		},
	},
	{
		fn: func() error {
			time.Sleep(time.Second * 1)
			results <- "Task3 completed"
			return nil
		},
	},
}

func TestNewPool(t *testing.T) {
	pool1 := NewPool(tasks, 10, 2)
	pool2 := NewPool(tasks, 10, -10)

	assert.Equal(t, pool1, pool2, "NewPool should return singleton struct")
	assert.Equal(t, 2, pool1.numCPUs, "No of CPUs used by pool should be 2")
}

func TestSetNumCPUs(t *testing.T) {
	pool := NewPool(tasks, 10, 4)

	assert.Equal(t, 4, pool.numCPUs, "No of CPUs used by pool should be 4")
	pool.SetNumCPUs(2)
	assert.Equal(t, 2, pool.numCPUs, "No of CPUs used by pool should be 2")
}

func TestPool_Run(t *testing.T) {
	pool := NewPool(tasks, 10, 4)
	pool.Run()
	assert.Equal(t, <-results, "Task3 completed", "Task3 should be completed first")
	assert.Equal(t, <-results, "Task2 completed", "Task2 should be completed second")
	assert.Equal(t, <-results, "Task1 completed", "Task1 should be completed last")
}
