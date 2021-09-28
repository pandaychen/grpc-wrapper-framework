package main

import (
	"fmt"
	"grpc-wrapper-framework/atreus"
	"grpc-wrapper-framework/config"
	pb "grpc-wrapper-framework/proto"

	"golang.org/x/net/context"

	//"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {

	s := &config.AtreusSvcConfig{
		Addr:                "127.0.0.1:50001",
		RegisterType:        "etcd",
		RegisterEndpoints:   "http://127.0.0.1:2379;",
		RegisterRootPath:    "/t",
		RegisterService:     "test",
		RegisterServiceVer:  "1.0",
		RegisterServiceAddr: "127.0.0.1:50001",
	}
	/*
		conn, err := grpc.Dial("127.0.0.1:50001", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
	*/
	conn := atreus.NewClient(s)

	client := pb.NewGreeterServiceClient(conn.RpcPersistClient)
	//add request
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(atreus.DefaultAtreusReqIDKey, "cvalue"))
	var fail int
	for i := 0; i < 1000; i++ {
		_, err := client.SayHello(ctx, &pb.HelloRequest{Name: "hello golang"})
		if err != nil {
			fmt.Println(err)
			fail++
		}
		//fmt.Printf("normal hello: resp=%v, error=%v\n", resp, err)
	}
	fmt.Println(fail)
}
