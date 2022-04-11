package main

import (
	"flag"
	"fmt"
	"grpc-wrapper-framework/atreus"
	"grpc-wrapper-framework/config"
	pb "grpc-wrapper-framework/proto"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
)

var (
	port = flag.Int("port", 12345, "listening port")
)

type xServer struct {
	BindAddr string
}

func main() {
	flag.Parse()

	BindAddr := fmt.Sprintf("127.0.0.1:%d", *port)
	lis, err := net.Listen("tcp", BindAddr)
	if err != nil {
		panic(err)
	}

	config.InitConfigAbsolutePath("../conf/", "grpc_server", "yaml")
	config.AtreusSvcConfigInit()

	grpclog.Infof("Server binding in %s...", BindAddr)
	s := atreus.NewServer(config.GetAtreusSvcConfig())
	pb.RegisterGreeterServiceServer(s.GetServer(), &xServer{})
	go s.Serve(lis)
	s.ExitWithSignalHandler()
}

func (xServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println(atreus.GetGlobalReqIDFromContext(ctx))

	//设置服务端超时时间大于配置时间
	time.Sleep(10 * time.Second)
	return &pb.HelloReply{Message: "Hello, " + req.Name}, nil
}
