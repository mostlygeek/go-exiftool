package exiftool

import (
	"errors"
	"sync"
)

// Pool creates multiple stay open exiftool instances and spreads the work
// across them with a simple round robin distribution.
type Pool struct {
	sync.Mutex
	stayopens []*Stayopen
	c         int
	l         int
	stopped   bool
}

func (p *Pool) Extract(filename string) (*Metadata, error) {
	if p.stopped {
		return nil, errors.New("Stopped")
	}
	p.Lock()
	p.c++
	key := p.c % p.l
	p.Unlock()
	return p.stayopens[key].Extract(filename)
}

func (p *Pool) Stop() {
	p.Lock()
	defer p.Unlock()
	for _, s := range p.stayopens {
		s.Stop()
	}
	p.stopped = true
}

func NewPool(exiftool string, num int) *Pool {
	p := &Pool{
		stayopens: make([]*Stayopen, num, num),
		l:         num,
	}

	for i := 0; i < num; i++ {
		p.stayopens[i] = NewStayopen(exiftool)
	}

	return p
}
