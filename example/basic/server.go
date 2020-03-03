package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/pandaychen/grpc-wrapper-framework/atreus"
	pb "github.com/pandaychen/grpc-wrapper-framework/proto"
	"golang.org/x/net/context"
)

var (
	port = flag.Int("port", 50001, "listening port")
)

type xServer struct {
	BindAddr string
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		panic(err)
	}
	s := atreus.NewServer()
	pb.RegisterGreeterServiceServer(s, &xServer{})
	s.Serve(lis)
}

func (xServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello, " + req.Name}, nil
}
