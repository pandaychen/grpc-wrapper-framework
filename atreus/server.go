package atreus

//A wrapper grpc-server

import (
	"google.golang.org/grpc"
)

func NewServer() *grpc.Server {
	var opt []grpc.ServerOption
	opt = append(opt, grpc.UnaryInterceptor(AtreusUnaryInterceptorChain(Recovery, Middle, Timing, Middle)))
	//return grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptorChain(Recovery, Logging)))
	return grpc.NewServer(opt...)
}
