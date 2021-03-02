package atreus

//A wrapper grpc-server

import (
	"sync"

	etcdv3 "go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	zaplog "github.com/pandaychen/goes-wrapper/zaplog"
	"github.com/pandaychen/grpc-wrapper-framework/common/vars"
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

	return grpc.NewServer(opt...)
}

// Server return the grpc server for registering service.
func (s *Server) Server() *grpc.Server {
	return s.RpcServer
}

// Use attachs a global inteceptor to the server
func (s *Server) Use(handlers ...grpc.UnaryServerInterceptor) *Server {
	new_size := len(s.InnerHandlers) + len(handlers)
	if new_size >= int(vars.ATREUS_MAX_INTERCEPTOR_NUM) {
		//限制拦截器的使用上限
		panic("too many interceptors")
	}
	mergedHandlers := make([]grpc.UnaryServerInterceptor, new_size)

	//warning: Should keep interceptors order
	copy(mergedHandlers, s.InnerHandlers)

	//copy new handles
	copy(mergedHandlers[len(s.InnerHandlers):], handlers)
	s.InnerHandlers = mergedHandlers
	return s
}
