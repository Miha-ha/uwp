package uwp

import (
	"sync"

	"github.com/pkg/errors"
)

type TaskFunc func() error

type Pool struct {
	concurency int
	tasksCh    chan TaskFunc
	wg         sync.WaitGroup
	err        error
}

func NewPool(concurency int) *Pool {
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
		p.wg.Add(1)
		err := task()
		if err != nil {
			if p.err == nil {
				p.err = err
			} else {
				p.err = errors.Wrapf(p.err, "%v", err)
			}
		}
		p.wg.Done()
	}
}

func (p *Pool) Add(fn TaskFunc) *Pool {
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
