package selectors

import (
	"container/ring"
	"sync"
)

type RoundRobin struct {
	r *ring.Ring
	l sync.RWMutex
}

func NewRoundRobin(strs []string) *RoundRobin {
	r := ring.New(len(strs))
	for _, s := range strs {
		r.Value = s
		r = r.Next()
	}
	return &RoundRobin{r: r}
}

func (rr *RoundRobin) Len() int {
	rr.l.RLock()
	defer rr.l.RUnlock()
	return rr.r.Len()
}

func (rr *RoundRobin) Choose() string {
	rr.l.Lock()
	defer rr.l.Unlock()
	if rr.r == nil {
		return ""
	}
	n := rr.r.Value.(string)
	rr.r = rr.r.Next()
	return n
}

func (rr *RoundRobin) Add(s string) {
	rr.l.Lock()
	defer rr.l.Unlock()
	nr := &ring.Ring{Value: s}
	if rr.r == nil {
		rr.r = nr
	} else {
		rr.r = rr.r.Link(nr).Next()
	}
}

func (rr *RoundRobin) Remove(s string) {
	rr.l.Lock()
	defer rr.l.Unlock()
	r := rr.r
	if rr.r.Len() == 1 {
		rr.r = ring.New(0)
		return
	}

	for i := rr.r.Len(); i > 0; i-- {
		r = r.Next()
		ba := r.Value.(string)
		if s == ba {
			rr.r = r.Unlink(1)
			return
		}
	}
}
