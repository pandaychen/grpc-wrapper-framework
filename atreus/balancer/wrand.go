package balancer

// balancer：带权重的随机算法

import (
	"context"
	mrand "math/rand"
	"strconv"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

func newRandomBuilder() balancer.Builder {
	return base.NewBalancerBuilderWithConfig(string(BALANCER_RandomWeight_NAME), &randomPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	//register randombuild to balancer
	balancer.Register(newRandomBuilder())
}

type randomPickerBuilder struct{}

func (*randomPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn

	for addr, sc := range readySCs {
		weight := 1
		m, ok := addr.Metadata.(*map[string]string)
		w, ok := (*m)["weight"]
		if ok {
			n, err := strconv.Atoi(w)
			if err == nil && n > 0 {
				weight = n
			}
		}
		for i := 0; i < weight; i++ {
			scs = append(scs, sc)
		}
	}

	// 初始化randomPicker
	return &randomPicker{
		subConns: scs,
	}
}

type randomPicker struct {
	subConns []balancer.SubConn
	lock     sync.Mutex
}

//Once Pick One available connection
func (p *randomPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	p.lock.Lock()
	sc := p.subConns[mrand.Intn(len(p.subConns))]
	p.lock.Unlock()
	return sc, nil, nil
}
