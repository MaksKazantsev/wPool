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

	wg sync.WaitGroup
}

var _ Pool = &pool{}

func NewPool(maxWorkers int) *pool {
	p := &pool{
		stopCh:     make(chan struct{}, 1),
		taskCh:     make(chan Tasker, maxWorkers),
		errorCh:    make(chan error, maxWorkers),
		maxWorkers: maxWorkers,
		wg:         sync.WaitGroup{},
	}

	p.run()

	return p
}

func (p *pool) run() {
	p.wg.Add(p.maxWorkers)

	for i := 0; i < p.maxWorkers; i++ {
		go func() {
			defer p.wg.Done()

			for {
				select {
				case t, ok := <-p.taskCh:
					if !ok {
						return
					}
					if err := t(); err != nil {
						p.errorCh <- err
					}
				case <-p.stopCh:
					return
				}
			}
		}()
	}
	go func() {
		p.wg.Wait()
		close(p.errorCh)
	}()
}

func (p *pool) Stop() {
	close(p.taskCh)
	close(p.stopCh)
}

func (p *pool) Task(tasker Tasker) {
	select {
	case p.taskCh <- tasker:
	}
}

func (p *pool) CatchError() chan error {
	return p.errorCh
}
