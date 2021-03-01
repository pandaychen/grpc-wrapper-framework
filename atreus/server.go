package atreus

//A wrapper grpc-server

import (
	"sync"

	etcdv3 "go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

//grpc-server核心结构（封装）
type Server struct {
	Logger *zap.Logger
	Conf   *AtreusServerConfig
	Lock   *sync.RWMutex

	//wrapper Server
	Server        *grpc.Server //原生Server
	EtcdClient    *etcdv3.Client
	InnerHandlers []grpc.UnaryServerInterceptor //拦截器数组
}

func NewServer() *grpc.Server {
	var opt []grpc.ServerOption
	opt = append(opt, grpc.UnaryInterceptor(AtreusUnaryInterceptorChain(Recovery, Middle, Timing, Middle)))
	//return grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptorChain(Recovery, Logging)))
	return grpc.NewServer(opt...)
}
