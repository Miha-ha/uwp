package uwp

import (
	"sync"

	"github.com/pkg/errors"
)

type TaskFunc func() error

type Pool struct {
	concurency int
	tasksCh    chan TaskFunc
	wgMu       sync.Mutex
	wg         sync.WaitGroup
	errMu      sync.Mutex
	err        error
}

func New(concurency int) *Pool {
	return &Pool{
		concurency: concurency,
		tasksCh:    make(chan TaskFunc),
	}
}

func (p *Pool) Run() *Pool {
	for i := 0; i < p.concurency; i++ {
		go p.work()
	}
	return p
}

func (p *Pool) work() {
	for task := range p.tasksCh {

		err := task()
		if err != nil {
			p.errMu.Lock()
			if p.err == nil {
				p.err = err
			} else {
				p.err = errors.Wrapf(p.err, "%v", err)
			}
			p.errMu.Unlock()
		}

		p.wgMu.Lock()
		p.wg.Done()
		p.wgMu.Unlock()
	}
}

func (p *Pool) Add(fn TaskFunc) *Pool {
	p.wg.Add(1)
	p.tasksCh <- fn
	return p
}

func (p *Pool) Wait() *Pool {
	p.wg.Wait()
	return p
}

func (p *Pool) Close() error {
	close(p.tasksCh)
	return p.err
}
