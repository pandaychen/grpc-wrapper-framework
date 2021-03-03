package atreus

//A wrapper grpc-server

import (
	"sync"
	"time"

	etcdv3 "go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	zaplog "github.com/pandaychen/goes-wrapper/zaplog"
	"github.com/pandaychen/grpc-wrapper-framework/config"
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

	IsDebug bool
}

func NewServer(conf *config.AtreusSvcConfig, opt ...grpc.ServerOption) *Server {
	if conf == nil {
		panic("atreus config null")
	}
	/*
		var opt []grpc.ServerOption
		opt = append(opt, grpc.UnaryInterceptor(AtreusUnaryInterceptorChain(Recovery, Middle, Timing, Middle)))
		//return grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptorChain(Recovery, Logging)))
	*/

	logger, _ := zaplog.ZapLoggerInit(DEFAULT_ATREUS_SERVICE_NAME)
	srv := &Server{
		Logger: logger,
		Lock:   new(sync.RWMutex),
		//InnerHandlers: make([]grpc.UnaryServerInterceptor, 0),
		Conf: NewAtreusServerConfig2(conf),
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

	opt = append(opt, keepaliveopts, grpc.UnaryInterceptor(srv.InnerHandlers))

	srv.RpcServer = grpc.NewServer(opt...)

	//Fill the interceptors

	srv.Use(s.Recovery())

	return srv
}

// Server return the grpc server for registering service.
func (s *Server) Server() *grpc.Server {
	return s.RpcServer
}
