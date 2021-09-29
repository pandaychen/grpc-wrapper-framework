package atreus

import (
	"fmt"
	"sync"

	"grpc-wrapper-framework/common/enums"
	"grpc-wrapper-framework/common/vars"
	"grpc-wrapper-framework/config"
	dis "grpc-wrapper-framework/microservice/discovery"
	com "grpc-wrapper-framework/microservice/discovery/common"

	zaplog "github.com/pandaychen/goes-wrapper/zaplog"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	Logger        *zap.Logger
	Conf          *AtreusClientConfig //指向 客户端配置
	Lock          *sync.RWMutex
	DialOpts      []grpc.DialOption             //grpc-客户端option
	InnerHandlers []grpc.UnaryClientInterceptor //GRPC拦截器数组

	RpcPersistClient *grpc.ClientConn

	CliResolver dis.ServiceResolverWrapper

	//sony breaker
	CbBreakerMap    map[string]*gobreaker.CircuitBreaker
	CbBreakerConfig gobreaker.Settings //这里暂时全局配置
}

func (c *Client) AddCliOpt(opts ...grpc.DialOption) *Client {
	c.DialOpts = append(c.DialOpts, opts...)
	return c
}

func NewClient(config *config.AtreusSvcConfig) *Client {
	var (
		err  error
		conn *grpc.ClientConn
	)

	logger, _ := zaplog.ZapLoggerInit(DEFAULT_ATREUS_SERVICE_NAME)

	cli := &Client{
		Logger:        logger,
		Lock:          new(sync.RWMutex),
		InnerHandlers: make([]grpc.UnaryClientInterceptor, 0),
		Conf:          NewAtreusClientConfig2(config),
		CbBreakerMap:  make(map[string]*gobreaker.CircuitBreaker),
	}

	cli.CliResolver, err = dis.NewDiscoveryResolver(&com.ResolverConfig{
		RegisterType:   enums.RegType(config.RegisterType),
		RootName:       config.RegisterRootPath,
		ServiceName:    config.RegisterService,
		ServiceVersion: config.RegisterServiceVer,
		Endpoint:       config.RegisterEndpoints,
		Schemename:     cli.Conf.DialScheme,
		Logger:         logger,
	})

	if err != nil {
		logger.Error("[NewClient]NewDiscoveryResolver error", zap.String("errmsg", err.Error()))
		panic(err)
	}

	//init client interceptors
	cli.Use(cli.Recovery(), cli.Timing(), cli.CircuitBreaker())

	//set dial options
	//TODO：配置化
	cli.DialOpts = append(cli.DialOpts, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithBalancerName("round_robin"), grpc.WithUnaryInterceptor(cli.BuildUnaryInterceptorChain2()))

	if cli.Conf.k8sSign && cli.Conf.DnsSign {
		//support K8S environment
		dial_address := fmt.Sprintf("dns:///%s:%d", cli.Conf.ServiceName, cli.Conf.ServicePort)
		conn, err = grpc.Dial(dial_address, cli.DialOpts...)
	} else {
		//TODO：FIX etcd config
		conn, err = grpc.Dial("etcdv3"+":///", cli.DialOpts...)
	}

	if err != nil {
		logger.Error("[NewClient]Dial Service error", zap.String("errmsg", err.Error()))
		panic(err)
	}

	//init breaker config
	cli.CbBreakerConfig.Name = ""
	cli.CbBreakerConfig.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	cli.RpcPersistClient = conn

	return cli
}

// Use方法为grpc的客户端添加一个全局拦截器，传入参数是多个grpc.UnaryClientInterceptor
func (c *Client) Use(handlers ...grpc.UnaryClientInterceptor) *Client {
	new_size := len(c.InnerHandlers) + len(handlers)
	if new_size >= int(vars.ATREUS_MAX_INTERCEPTOR_NUM) {
		//拦截器过多影响处理性能和延迟
		panic("too many client handlers")
	}

	//将参数中的handlers添加在已有拦截器序列的后面，经典的复制slice的方法
	mergedHandlers := make([]grpc.UnaryClientInterceptor, new_size)
	copy(mergedHandlers, c.InnerHandlers)
	copy(mergedHandlers[len(c.InnerHandlers):], handlers)

	//new interceptors
	c.InnerHandlers = mergedHandlers
	return c
}

// 实现链式的客户端拦截器
func (c *Client) BuildUnaryInterceptorChain2() grpc.UnaryClientInterceptor {
	var (
		size int = len(c.InnerHandlers)
	)

	if size == 0 {
		//return grpc.UnaryClientInterceptor{}
		return func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}

	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var (
			i            int
			chainHandler grpc.UnaryInvoker
		)
		chainHandler = func(ictx context.Context, imethod string, ireq, ireply interface{}, ic *grpc.ClientConn, iopts ...grpc.CallOption) error {
			if i == size-1 {
				//返回RPC调用
				return invoker(ictx, imethod, ireq, ireply, ic, iopts...)
			}
			i++
			return c.InnerHandlers[i](ictx, imethod, ireq, ireply, ic, chainHandler, iopts...)
		}

		//返回第0号位置上的烂机器
		return c.InnerHandlers[0](ctx, method, req, reply, cc, chainHandler, opts...)
	}
}
