package balancer

//按照P2C的思路，选择最小连接的节点

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"

	"github.com/pandaychen/goes-wrapper/pymath"
)

type leastConnNode struct {
	balancer.SubConn       //backend
	inflight         int64 //当前活跃（请求中）连接计数
}

type LeastConnPicker struct {
	sync.Mutex
	nodes  []*leastConnNode
	rand   *rand.Rand
	logger *zap.Logger
}

type leastConnPickerBuilder struct {
	logger *zap.Logger
}

func newLeastConnBuilder(logger *zap.Logger) balancer.Builder {
	return base.NewBalancerBuilder(string(BALANCER_LeastConn_NAME), &leastConnPickerBuilder{
		logger: logger,
	}, base.Config{HealthCheck: true})
}

func RegisterLeastConnPickerBuilder(logger *zap.Logger) {
	balancer.Register(newLeastConnBuilder(logger))
}

// 每当有后端节点上下线时触发，balance pool重新生成
func (b *leastConnPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.Picker {
	var (
		nodes []*leastConnNode
	)
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	//init node，应该去掉已经存在的节点
	for subConn, _ := range buildInfo.ReadySCs {
		nodes = append(nodes, &leastConnNode{subConn, 0})
	}

	return &LeastConnPicker{
		logger: b.logger,
		nodes:  nodes,
		rand:   rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (p *LeastConnPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var (
		indexA, indexB int
		pickResult     balancer.PickResult
		pickNode       *leastConnNode
	)

	if len(p.nodes) == 0 {
		// 无可用
		return pickResult, balancer.ErrNoSubConnAvailable
	} else if len(p.nodes) == 1 {
		pickNode = p.nodes[0]
	} else {
		p.Lock()
		indexA, indexB = pymath.PowerOfTwoChoices(p.rand, len(p.nodes))
		p.Unlock()

		//选择连接数较小的node
		if p.nodes[indexA].inflight < p.nodes[indexB].inflight {
			pickNode = p.nodes[indexA]
		} else {
			pickNode = p.nodes[indexB]
		}
	}

	atomic.AddInt64(&pickNode.inflight, 1)

	pickResult.SubConn = pickNode
	pickResult.Done = func(info balancer.DoneInfo) {
		// RPC请求结束，连接减1
		atomic.AddInt64(&pickNode.inflight, -1)
	}

	return pickResult, nil
}
