package balancer

import (
	"github.com/pkg/errors"

	//"math/rand"
	//"strconv"

	"sync"

	pybalancer "github.com/pandaychen/goes-wrapper/balancer"

	"go.uber.org/zap"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	//"google.golang.org/grpc/resolver"
)

type simpleRoundRobinPicker struct {
	subConns []balancer.SubConn //一个balancer.SubConn标识一个长连接,subConns 标识所有活动连接数组
	lock     sync.Mutex
	Index    int
	Wrr      *pybalancer.NginxWeightRoundrobin
	Logger   *zap.Logger
}

// newsimpleRoundRobinBuilder creates a new roundrobin balancer builder
func newsimpleRoundRobinBuilder(logger *zap.Logger) balancer.Builder {
	return base.NewBalancerBuilder(string(BALANCER_SimpleWeightRR_NAME),
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
func (r *simpleRoundRobinPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.Picker {
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	wrr := pybalancer.NewNginxWeightRoundrobin(r.Logger)
	var scs []balancer.SubConn
	for subconn, sc := range buildInfo.ReadySCs {
		weight := GetServerWeightValue(sc.Address)
		// 将conn存储在WRR的slice中
		wrr.AddNode(sc.Address.Addr, subconn, weight)
		scs = append(scs, subconn)
	}
	//Build的作用是：根据readyScs，构造LB算法选择用的初始化集合，当然可以根据权重对subConns进行调整
	return &simpleRoundRobinPicker{
		subConns: scs,
		Wrr:      wrr,
		Logger:   r.Logger,
	}
}

//Picker方法：每次客户端RPC-CALL都会调用
func (p *simpleRoundRobinPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	var (
		pickResult balancer.PickResult
	)
	p.lock.Lock()
	defer p.lock.Unlock()
	sc := p.Wrr.GetNextNode()

	if sc == nil {
		return pickResult, errors.New("Pick one connection error")
	}

	if _, ok := sc.NodeMetadata.(balancer.SubConn); !ok {
		return pickResult, errors.New("system error")
	}

	pickResult.SubConn = sc.NodeMetadata.(balancer.SubConn)

	return pickResult, nil
}
