package xmath

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type XProbability struct {
	// rand.New(...) returns a non thread safe object
	r *rand.Rand
	sync.Mutex
	sampling float64 //0 ~ 1
}

func NewProbability(sampling float64) *XProbability {
	return &XProbability{
		sampling: sampling,
		r:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (p *XProbability) TrueOrNot() bool {
	p.Lock()
	t := p.r.Float64() < p.sampling
	p.Unlock()
	return t
}

func (p *XProbability) TrueOrNotWithProbable(proba float64) bool {
	p.Lock()
	t := p.r.Float64() < proba
	p.Unlock()
	return t
}

func main() {
	p := NewProbability(0.3)
	fmt.Println(p.TrueOrNotWithProbable(0.3))
}
