package wPool

import "sync"

type Tasker func() error

type Pool interface {
	Task(tasker Tasker)
	CatchError() chan error
	Stop()
}

type pool struct {
	maxWorkers int

	errorCh chan error
	taskCh  chan Tasker
	stopCh  chan struct{}
}

var _ Pool = &pool{}

func NewPool(maxWorkers int) *pool {
	p := &pool{
		stopCh:     make(chan struct{}, 1),
		taskCh:     make(chan Tasker, maxWorkers),
		errorCh:    make(chan error, maxWorkers),
		maxWorkers: maxWorkers,
	}

	p.run()

	return p
}

func (p *pool) run() {
	wg := sync.WaitGroup{}
	wg.Add(p.maxWorkers)

	for range p.maxWorkers {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			for t := range p.taskCh {
				if err := t(); err != nil {
					p.errorCh <- err
				}
			}
		}(&wg)
	}
	go func() {
		wg.Wait()
		close(p.errorCh)
	}()
}

func (p *pool) Stop() {
	close(p.taskCh)
}

func (p *pool) Task(tasker Tasker) {
	select {
	case p.taskCh <- tasker:
	}
}

func (p *pool) CatchError() chan error {
	return p.errorCh
}
