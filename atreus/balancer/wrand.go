package balancer

// balancer：带权重的随机算法

import (
	mrand "math/rand"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

func newRandomBuilder(logger *zap.Logger) balancer.Builder {
	return base.NewBalancerBuilder(string(BALANCER_RandomWeight_NAME), &randomPickerBuilder{
		Logger: logger,
	}, base.Config{HealthCheck: true})
}

func RegisterRandomBuilderPickerBuilder(logger *zap.Logger) {
	//register randombuild to balancer
	balancer.Register(newRandomBuilder(logger))
}

type randomPickerBuilder struct {
	Logger *zap.Logger
}

func (*randomPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.Picker {
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var (
		scs    []balancer.SubConn
		weight int
	)

	for subconn, sc := range buildInfo.ReadySCs {
		//get node weight
		weight = GetServerWeightValue(sc.Address.Metadata)
		for i := 0; i < weight; i++ {
			scs = append(scs, subconn)
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
func (p *randomPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	var (
		pickResult balancer.PickResult
	)
	p.lock.Lock()
	pickResult.SubConn = p.subConns[mrand.Intn(len(p.subConns))]
	p.lock.Unlock()
	return pickResult, nil
}
