package Go_Pool

import (
	"errors"
	"runtime"
	"sync"
)

type Pool struct {
	concurrency      int
	numCPUs          int
	tasksChan        chan *Task
	activeGoRoutines chan bool
	wg               sync.WaitGroup
}

var pool *Pool
var once sync.Once

func NewPool(concurrency int, numCPUs int) *Pool {
	once.Do(func() {
		pool = &Pool{
			concurrency:      concurrency,
			tasksChan:        make(chan *Task, concurrency),
			activeGoRoutines: make(chan bool, concurrency),
		}

		defer func() { go pool.schedule() }()
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

func (p *Pool) work() {
	for task := range p.tasksChan {
		task.Run(&p.wg)
	}
}

func (p *Pool) AddTask(task *Task) {
	defer p.wg.Wait()

	p.wg.Add(1)
	p.tasksChan <- task

}

func (p *Pool) schedule() {
	for task := range p.tasksChan {
		p.activeGoRoutines <- true
		go func() {
			defer func() {
				<-p.activeGoRoutines
			}()
			task.Run(&p.wg)
		}()
	}
}

func (p *Pool) Close() {
	close(p.tasksChan)
	close(p.activeGoRoutines)
}
