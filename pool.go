package exiftool

// Pool creates multiple stay open exiftool instances and spreads the work
// across them with a simple round robin distribution.
type Pool struct {
	stayopens []*Stayopen
	c         uint8
	l         int
}

func (p *Pool) Extract(filename string) (*Metadata, error) {
	p.c++
	return p.stayopens[p.c%l].Extract(filename)
}

func (p *Pool) Stop() {
	for _, s := range p.stayopens {
		s.Stop()
	}
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
