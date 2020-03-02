package atreus

//A wrapper grpc-server

import (
	"google.golang.org/grpc"
)

func NewServer() *grpc.Server {
	var opt []grpc.ServerOption
	return grpc.NewServer(opt...)
}
