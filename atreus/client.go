package atreus

import (
	"fmt"

	zaplog "github.com/pandaychen/goes-wrapper/zaplog"
	"github.com/pandaychen/grpc-wrapper-framework/common/enums"
	"github.com/pandaychen/grpc-wrapper-framework/config"
	dis "github.com/pandaychen/grpc-wrapper-framework/microservice/discovery"
	com "github.com/pandaychen/grpc-wrapper-framework/microservice/discovery/common"
	"google.golang.org/grpc"
)

type Client struct {
	Conf *AtreusClientConfig //指向 客户端配置

	DialOpts []grpc.DialOption             //grpc-客户端option
	handlers []grpc.UnaryClientInterceptor //GRPC拦截器数组

	RpcPersistClient *grpc.ClientConn
}

func (c *Client) AddCliOpt(opts ...grpc.DialOption) *Client {
	c.DialOpts = append(c.DialOpts, opts...)
	return c
}

func NewClient(conf *config.AtreusSvcConfig) *Client {
	logger, _ := zaplog.ZapLoggerInit(DEFAULT_ATREUS_SERVICE_NAME)
	_, err := dis.NewDiscoveryResolver(&com.ResolverConfig{
		RegisterType:   enums.RegType(conf.RegisterType),
		RootName:       conf.RegisterRootPath,
		ServiceName:    conf.RegisterService,
		ServiceVersion: conf.RegisterServiceVer,
		Endpoint:       conf.RegisterEndpoints,
		Schemename:     "etcdv3",
		Logger:         logger,
	})
	conn, err := grpc.Dial("etcdv3"+":///", grpc.WithBalancerName("round_robin"), grpc.WithBlock(), grpc.WithInsecure())
	fmt.Println(conn, err)
	client := &Client{
		RpcPersistClient: conn,
	}
	return client
}
