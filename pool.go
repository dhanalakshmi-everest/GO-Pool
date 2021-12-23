package Go_Pool

import (
	"fmt"
	"runtime"
	"sync"
)

type Pool struct {
	Tasks       []*Task
	concurrency int
	numCPUs     int
	tasksChan   chan *Task
	wg          sync.WaitGroup
}

var pool *Pool
var once sync.Once

func NewPool(tasks []*Task, concurrency int, numCPUs int) *Pool {
	once.Do(func() {
		fmt.Println()
		pool = &Pool{
			Tasks:       tasks,
			concurrency: concurrency,
			tasksChan:   make(chan *Task),
		}
	})
	SetNumCPUs(numCPUs)

	return pool
}

func SetNumCPUs(numCPUs int) {
	pool.numCPUs = runtime.GOMAXPROCS(numCPUs)
}

func (p *Pool) Run() {
	for i := 0; i < p.concurrency; i++ {
		go p.work()
	}

	p.wg.Add(len(p.Tasks))
	for _, task := range p.Tasks {
		p.tasksChan <- task
	}

	close(p.tasksChan)

	p.wg.Wait()
}

func (p *Pool) work() {
	for task := range p.tasksChan {
		task.Run(&p.wg)
	}
}
