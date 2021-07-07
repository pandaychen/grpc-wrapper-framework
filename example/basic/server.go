package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/pandaychen/grpc-wrapper-framework/atreus"
	"github.com/pandaychen/grpc-wrapper-framework/config"
	pb "github.com/pandaychen/grpc-wrapper-framework/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
)

var (
	port = flag.Int("port", 50003, "listening port")
)

type xServer struct {
	BindAddr string
}

func main() {
	flag.Parse()

	//init logger
	lc := config.LogConfig{}
	grpclog.SetLoggerV2(lc.CreateNewLogger("grpc-basic-service"))

	BindAddr := fmt.Sprintf("0.0.0.0:%d", *port)
	lis, err := net.Listen("tcp", BindAddr)
	if err != nil {
		panic(err)
	}

	grpclog.Infof("Server binding in %s...", BindAddr)
	s := atreus.NewServer(&config.AtreusSvcConfig{
		Addr:                "127.0.0.1:50003",
		RegisterType:        "etcd",
		RegisterEndpoints:   "http://127.0.0.1:2379;",
		RegisterRootPath:    "/t",
		RegisterService:     "test",
		RegisterServiceVer:  "1.0",
		RegisterServiceAddr: "127.0.0.1:50003",
	})
	pb.RegisterGreeterServiceServer(s.GetServer(), &xServer{})
	s.Serve(lis)
}

func (xServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println(atreus.GetGlobalReqIDFromContext(ctx))
	return &pb.HelloReply{Message: "Hello, " + req.Name}, nil
}
