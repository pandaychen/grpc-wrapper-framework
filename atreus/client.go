package atreus

import (
	"fmt"
	"sync"

	zaplog "github.com/pandaychen/goes-wrapper/zaplog"
	"github.com/pandaychen/grpc-wrapper-framework/common/enums"
	"github.com/pandaychen/grpc-wrapper-framework/config"
	dis "github.com/pandaychen/grpc-wrapper-framework/microservice/discovery"
	com "github.com/pandaychen/grpc-wrapper-framework/microservice/discovery/common"
	"go.uber.org/zap"
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

	if cli.Conf.k8sSign && cli.Conf.DnsSign {
		//support TKE environment
		dial_address := fmt.Sprintf("dns:///%s:%d", cli.Conf.ServiceName, cli.Conf.ServicePort)
		conn, err = grpc.Dial(dial_address, grpc.WithBalancerName("round_robin"), grpc.WithBlock(), grpc.WithInsecure())
	} else {
		conn, err = grpc.Dial("etcdv3"+":///", grpc.WithBalancerName("round_robin"), grpc.WithBlock(), grpc.WithInsecure())
	}

	if err != nil {
		logger.Error("[NewClient]Dial Service error", zap.String("errmsg", err.Error()))
		panic(err)
	}

	cli.RpcPersistClient = conn

	return cli
}
