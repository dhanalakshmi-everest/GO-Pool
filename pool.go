package Go_Pool

import (
	"errors"
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
		pool = &Pool{
			Tasks:       tasks,
			concurrency: concurrency,
			tasksChan:   make(chan *Task, concurrency),
		}
	})
	pool.SetNumCPUs(numCPUs)

	return pool
}

func (p *Pool) SetNumCPUs(numCPUs int) error {
	if numCPUs < 1 {
		return errors.New("No of CPUs is a negative number")
	}

	runtime.GOMAXPROCS(numCPUs)
	p.numCPUs = numCPUs

	return nil
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
