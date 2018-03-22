package exiftool

import (
	"math/rand"
	"time"
)

type Pool struct {
	stayopens []*Stayopen
	rand      *rand.Rand
}

func (p *Pool) Extract(filename string) (*Metadata, error) {
	return p.stayopens[p.rand.Intn(len(p.stayopens))].Extract(filename)
}

func (p *Pool) Stop() {
	for _, s := range p.stayopens {
		s.Stop()
	}
}

func NewPool(exiftool string, num int) *Pool {
	p := &Pool{
		stayopens: make([]*Stayopen, num, num),
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for i := 0; i < num; i++ {
		p.stayopens[i] = NewStayopen(exiftool)
	}

	return p
}
