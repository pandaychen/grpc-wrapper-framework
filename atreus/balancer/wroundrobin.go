package balancer

import (
	"context"
	"errors"

	//"math/rand"
	//"strconv"

	"sync"

	pybalancer "github.com/pandaychen/goes-wrapper/balancer"

	"go.uber.org/zap"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

type simpleRoundRobinPicker struct {
	subConns []balancer.SubConn //一个balancer.SubConn标识一个长连接,subConns 标识所有活动连接数组
	lock     sync.Mutex
	Index    int
	Wrr      *pybalancer.NginxWeightRoundrobin
}

// newsimpleRoundRobinBuilder creates a new roundrobin balancer builder
func newsimpleRoundRobinBuilder(logger *zap.Logger) balancer.Builder {
	return base.NewBalancerBuilderWithConfig(string(BALANCER_SimpleWeightRR_NAME),
		&simpleRoundRobinPickerBuilder{
			Logger: logger,
		},
		base.Config{HealthCheck: true})
}

/*
func init() {
	balancer.Register(newsimpleRoundRobinBuilder())
}
*/

func RegisterSimpleRoundRobinPickerBuilder(logger *zap.Logger) {
	balancer.Register(newsimpleRoundRobinBuilder(logger))
}

type simpleRoundRobinPickerBuilder struct {
	Logger *zap.Logger
}

// Triggers where grpc-client-pool changing(Once backend node On/Off)
func (r *simpleRoundRobinPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	r.Logger.Info("ready", zap.Any("scs", len(readySCs)))

	wrr := pybalancer.NewNginxWeightRoundrobin(r.Logger)

	var scs []balancer.SubConn
	for addr, sc := range readySCs {
		weight := GetServerWeightValue(addr.Metadata)
		wrr.AddNode(sc, weight)
	}
	//Build的作用是：根据readyScs，构造LB算法选择用的初始化集合，当然可以根据权重对subConns进行调整
	return &simpleRoundRobinPicker{
		subConns: scs,
		Wrr:      wrr,
	}
}

//Picker方法：每次客户端RPC-CALL都会调用
func (p *simpleRoundRobinPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	sc := p.Wrr.GetNextNode()

	if sc == nil {
		return nil, nil, errors.New("Pick one connection error")
	}

	if _, ok := sc.NodeMetadata.(balancer.SubConn); !ok {
		return nil, nil, errors.New("system error")
	}

	return sc.NodeMetadata.(balancer.SubConn), nil, nil
}
