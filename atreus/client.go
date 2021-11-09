package atreus

import (
	"errors"
	"fmt"
	"sync"

	"grpc-wrapper-framework/common/enums"
	"grpc-wrapper-framework/common/vars"
	"grpc-wrapper-framework/config"
	"grpc-wrapper-framework/logger"
	dis "grpc-wrapper-framework/microservice/discovery"
	com "grpc-wrapper-framework/microservice/discovery/common"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// 客户端封装结构
type Client struct {
	Logger        *zap.Logger
	Conf          *config.AtreusCliConfig //客户端配置
	Lock          *sync.RWMutex
	DialOpts      []grpc.DialOption             //grpc-客户端option
	InnerHandlers []grpc.UnaryClientInterceptor //GRPC拦截器数组

	CliResolver dis.ServiceResolverWrapper

	//sony breaker
	CbBreakerMap    map[string]*gobreaker.CircuitBreaker
	CbBreakerConfig gobreaker.Settings //这里暂时全局配置

	RpcPersistClient *grpc.ClientConn
}

func (c *Client) AddCliOpt(opts ...grpc.DialOption) *Client {
	c.DialOpts = append(c.DialOpts, opts...)
	return c
}

func NewClient(config *config.AtreusCliConfig) (*Client, error) {
	var (
		err       error
		conn      *grpc.ClientConn
		is_direct bool
	)

	logconf := logger.LogConfig{
		ServiceName: DEFAULT_ATREUS_SERVICE_NAME,
	}

	logger, err := logconf.CreateNewLogger(config.LogConf)
	if err != nil {
		return nil, err
	}
	cli := &Client{
		Logger:        logger,
		Lock:          new(sync.RWMutex),
		InnerHandlers: make([]grpc.UnaryClientInterceptor, 0),
		Conf:          config,
		CbBreakerMap:  make(map[string]*gobreaker.CircuitBreaker),
	}

	switch config.CliConf.DialScheme {
	case string(enums.RET_TYPE_DIRECT):
		is_direct = true
		cli.DialOpts = append(cli.DialOpts, grpc.WithBlock())
	case string(enums.REG_TYPE_DNS):
		cli.DialOpts = append(cli.DialOpts, grpc.WithBlock())
	case string(enums.REG_TYPE_ETCD):
		cli.CliResolver, err = dis.NewDiscoveryResolver(&com.ResolverConfig{
			RegisterType:   enums.RegType(config.RegistryConf.RegisterType),
			RootName:       config.RegistryConf.RegisterRootPath,
			ServiceName:    config.RegistryConf.RegisterService,
			ServiceVersion: config.RegistryConf.RegisterServiceVer,
			Endpoint:       config.RegistryConf.RegisterEndpoints,
			Schemename:     config.CliConf.DialScheme,
			Logger:         logger,
		})
	case string(enums.REG_TYPE_CONSUL):
		return nil, errors.New("not support dial scheme")
	default:
		return nil, errors.New("not support dial scheme")
	}

	if !is_direct {
		switch config.CliConf.LbType {
		case string(enums.LB_TYPE_RR):
			cli.DialOpts = append(cli.DialOpts, grpc.WithBalancerName("round_robin"))
		default:
			return nil, errors.New("not support lb type")
		}
	}

	if err != nil {
		logger.Error("[NewClient]NewDiscoveryResolver error", zap.String("errmsg", err.Error()))
		panic(err)
	}

	//init client interceptors
	cli.Use(cli.Recovery(), cli.Timing())

	//set dial options
	cli.DialOpts = append(cli.DialOpts, grpc.WithInsecure(), grpc.WithUnaryInterceptor(cli.BuildUnaryInterceptorChain2()))

	if config.BreakerConf.On {
		//init breaker config
		cli.CbBreakerConfig.Name = ""
		cli.CbBreakerConfig.ReadyToTrip = func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		}
		//add breaker
		cli.Use(cli.CircuitBreaker())
	}

	switch config.CliConf.DialScheme {
	case string(enums.RET_TYPE_DIRECT):
		dial_address := fmt.Sprintf("%s:%d", config.CliConf.DialAddress, config.CliConf.DialPort)
		conn, err = grpc.Dial(dial_address, cli.DialOpts...)
	case string(enums.REG_TYPE_DNS):
		//support K8S environment
		dial_address := fmt.Sprintf("dns:///%s:%d", config.SrvDnsConf.SrvName, config.SrvDnsConf.SrvPort)
		conn, err = grpc.Dial(dial_address, cli.DialOpts...)
	case string(enums.REG_TYPE_ETCD):
		conn, err = grpc.Dial(fmt.Sprintf("%s:///", cli.CliResolver.Scheme()), cli.DialOpts...)
	case string(enums.REG_TYPE_CONSUL):
		return nil, errors.New("not support dial scheme")
	default:
		return nil, errors.New("not support dial scheme")
	}

	if err != nil {
		logger.Error("[NewClient]Dial Service error", zap.String("errmsg", err.Error()))
		panic(err)
	}

	cli.RpcPersistClient = conn

	return cli, nil
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
