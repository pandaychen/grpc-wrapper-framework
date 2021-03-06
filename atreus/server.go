package atreus

//A wrapper grpc-server

import (
	"fmt"
	"net"
	"sync"
	"time"

	etcdv3 "go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"

	zaplog "github.com/pandaychen/goes-wrapper/zaplog"
	"github.com/pandaychen/grpc-wrapper-framework/common/enums"
	"github.com/pandaychen/grpc-wrapper-framework/common/vars"
	"github.com/pandaychen/grpc-wrapper-framework/config"
	com "github.com/pandaychen/grpc-wrapper-framework/microservice/discovery/common"
	"github.com/pandaychen/grpc-wrapper-framework/pkg/xrand"

	dis "github.com/pandaychen/grpc-wrapper-framework/microservice/discovery"
)

const (
	DEFAULT_ATREUS_SERVICE_NAME = "atreus_svc"
)

//grpc-server核心结构（封装）
type Server struct {
	Logger *zap.Logger
	Conf   *AtreusServerConfig
	Lock   *sync.RWMutex

	//wrapper Server
	RpcServer     *grpc.Server //原生Server
	EtcdClient    *etcdv3.Client
	InnerHandlers []grpc.UnaryServerInterceptor //拦截器数组
	ServiceReg    dis.ServiceRegisterWrapper
	IsDebug       bool
}

func NewServer(conf *config.AtreusSvcConfig, opt ...grpc.ServerOption) *Server {
	var err error
	if conf == nil {
		panic("atreus server config null")
	}
	/*
		var opt []grpc.ServerOption
		opt = append(opt, grpc.UnaryInterceptor(AtreusUnaryInterceptorChain(Recovery, Middle, Timing, Middle)))
		//return grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptorChain(Recovery, Logging)))
	*/

	logger, _ := zaplog.ZapLoggerInit(DEFAULT_ATREUS_SERVICE_NAME)
	srv := &Server{
		Logger:        logger,
		Lock:          new(sync.RWMutex),
		InnerHandlers: make([]grpc.UnaryServerInterceptor, 0),
		Conf:          NewAtreusServerConfig2(conf),
	}

	//初始化gRPC-Server的keepalive参数
	keepaliveopts := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: time.Duration(srv.Conf.IdleTimeout), //如果一个client空闲超过MaxConnectionIdle-s,发送一个GOAWAY,为了防止同一时间发送大量GOAWAY
		//假设MaxConnectionIdle=15s，那么会在15s时间间隔上下浮动MaxConnectionIdle*10%,即15+1.5或者15-1.5
		MaxConnectionAgeGrace: time.Duration(srv.Conf.ForceCloseWait),    //在强制关闭连接之间,允许有ForceCloseWait-s的时间完成pending的rpc请求
		Time:                  time.Duration(srv.Conf.KeepAliveInterval), //如果一个clinet空闲超过KeepAliveInterval-s,则发送一个ping请求
		Timeout:               time.Duration(srv.Conf.KeepAliveTimeout),  //如果ping请求KeepAliveTimeout-s内未收到回复,则认为该连接已断开
		MaxConnectionAge:      time.Duration(srv.Conf.MaxLifeTime),       //如果任意连接存活时间超过MaxLifeTime-s,发送一个GOAWAY
	})

	opt = append(opt, keepaliveopts, grpc.UnaryInterceptor(srv.BuildUnaryInterceptorChain2))

	srv.RpcServer = grpc.NewServer(opt...)

	//Fill the interceptors
	srv.Use(srv.Recovery(), srv.AtreusXRequestId(), srv.Metrics2Prometheus())

	nodeinfo := com.ServiceBasicInfo{
		AddressInfo: conf.Addr,
		Metadata:    metadata.Pairs(vars.SERVICE_WEIGHT_KEY, conf.InitWeight),
	}

	srv.ServiceReg, err = dis.NewDiscoveryRegister(&com.RegisterConfig{
		RegisterType:   enums.RegType(conf.RegisterType),
		RootName:       conf.RegisterRootPath,
		ServiceName:    conf.RegisterService,
		ServiceVersion: conf.RegisterServiceVer,
		ServiceNodeID:  fmt.Sprintf("addr#%s", conf.Addr),
		RandomSuffix:   string(xrand.RandomString(8)),
		Ttl:            conf.RegisterTTL,
		Endpoint:       conf.RegisterEndpoints,
		Logger:         logger,
		NodeData:       nodeinfo,
	})

	if err == nil {
		err = srv.ServiceReg.ServiceRegister()
		if err != nil {
			logger.Error("[NewServer]ServiceRegister error", zap.String("errmsg", err.Error()))
			return nil
		}
	} else {
		logger.Error("[NewServer]NewDiscoveryRegister error", zap.String("errmsg", err.Error()))
		panic(err)
	}

	return srv
}

// Server return the grpc server for registering service.
func (s *Server) GetServer() *grpc.Server {
	return s.RpcServer
}

func (s *Server) Serve(lis net.Listener) error {
	return s.RpcServer.Serve(lis)
}
