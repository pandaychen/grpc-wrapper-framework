package main

import (
	"flag"
	"fmt"

	"grpc-wrapper-framework/atreus"
	"grpc-wrapper-framework/config"
	pb "grpc-wrapper-framework/proto"

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

	//init logger
	//lc := config.LogConfig{}
	//grpclog.SetLoggerV2(lc.CreateNewLogger("grpc-basic-service"))

	BindAddr := fmt.Sprintf("127.0.0.1:%d", *port)

	config.InitConfigAbsolutePath("./", "server", "yaml")
	config.AtreusSvcConfigInit()
	fmt.Println(config.GetAtreusSvcConfig())

	grpclog.Infof("Server binding in %s...", BindAddr)
	s := atreus.NewServer(config.GetAtreusSvcConfig())
	pb.RegisterGreeterServiceServer(s.GetServer(), &xServer{})
	g, _ := atreus.NewGracefulGrpcAppserver(s, BindAddr)
	g.RunServer()
}

func (xServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println(atreus.GetGlobalReqIDFromContext(ctx))
	return &pb.HelloReply{Message: "Hello, " + req.Name}, nil
}
