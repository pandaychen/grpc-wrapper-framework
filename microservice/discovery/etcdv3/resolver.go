package etcdv3

import (
	"strings"
	"sync"

	com "grpc-wrapper-framework/microservice/discovery/common"

	etcd3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/resolver"
)

// Default shchme
const (
	defaultEtcdScheme = "etcdv3"
)

var Once sync.Once

type EtcdResolver struct {
	Schemename string
	EtcdCli    *etcd3.Client // etcd3 client
	WatchKey   string
	Watcher    *EtcdWatcher
	Clientconn resolver.ClientConn
	Wg         sync.WaitGroup
	//CloseCh    chan struct{} // 关闭 channel

	//control
	Ctx    *context.Context
	Cancel context.CancelFunc

	Logger *zap.Logger
}

func NewResolverRegister(config *com.ResolverConfig) (*EtcdResolver, error) {
	etcdConfg := etcd3.Config{
		Endpoints: strings.Split(config.Endpoint, ";"),
	}

	client, err := etcd3.New(etcdConfg)
	if err != nil {
		config.Logger.Error("[NewResolverRegister]Create etcdv3 client error", zap.String("errmsg", err.Error()))
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	r := &EtcdResolver{
		Schemename: config.Schemename,
		EtcdCli:    client,
		WatchKey:   config.BuildEtcdPrefix(),
		Logger:     config.Logger,
		Ctx:        &ctx,
		Cancel:     cancel,
	}

	if r.Schemename == "" {
		r.Schemename = defaultEtcdScheme
	}

	// 调用grpc/resolver包的全局方法注册
	Once.Do(
		func() {
			resolver.Register(r)
		})
	return r, nil
}

// Build returns itself for resolver, because it's both a builder and a resolver.
func (r *EtcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	//用来从etcd获取serverlist,并通知Clientconn更新连接池
	r.Clientconn = cc

	//create watcher
	r.Watcher = NewEtcdWatcher(r.Ctx, r.WatchKey, r.EtcdCli, r.Logger, &r.Wg)

	//start watcher，从etcd中监控最新的地址变化，并通知clientconn
	r.start()

	return r, nil
}

// Scheme returns the scheme.
func (r *EtcdResolver) Scheme() string {
	return r.Schemename
}

// ResolveNow is a noop for resolver.
func (r *EtcdResolver) ResolveNow(o resolver.ResolveNowOptions) {
}

// Close is a noop for resolver.
func (r *EtcdResolver) Close() {
	r.Watcher.Close()
	r.Wg.Wait()
}

// Start Resover return a closeCh, Should call by Builde func()
func (r *EtcdResolver) start() {
	addrlist_channel := r.Watcher.start()
	r.Wg.Add(1)
	go func() {
		defer r.Wg.Done()
		for addr := range addrlist_channel {
			//range在channel上,addr为最新的[]resolver.Address
			r.Clientconn.UpdateState(resolver.State{Addresses: addr})
		}
	}()
}
