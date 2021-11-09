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
	config.InitConfigAbsolutePath("./", "client", "yaml")
	config.AtreusCliConfigInit()
	/*
		conn, err := grpc.Dial("127.0.0.1:50001", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
	*/
	conn, _ := atreus.NewClient(config.GetAtreusCliConfig())

	client := pb.NewGreeterServiceClient(conn.RpcPersistClient)
	//add request
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(atreus.DefaultAtreusReqIDKey, "cvalue", "app", "test", "token", "test", "method", "normal"))
	var fail int
	for i := 0; i < 10; i++ {
		resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "hello golang"})
		if err != nil {
			fmt.Println(err)
			fail++
		} else {
			fmt.Printf("normal hello: resp=%v, error=%v\n", resp, err)
		}
	}
	fmt.Println(fail)
}
