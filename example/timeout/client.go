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

func extractIncomingAndClone(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.MD{}
	}

	return md.Copy()
}

func main() {
	var (
		md   metadata.MD // == nil
		fail int
	)
	config.InitConfigAbsolutePath("../conf", "grpc_client", "yaml")
	config.AtreusCliConfigInit()

	conn, _ := atreus.NewClient(config.GetAtreusCliConfig())
	md = metadata.MD{}
	md["sendParamA"] = nil
	client := pb.NewGreeterServiceClient(conn.RpcPersistClient)
	//add request
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(atreus.DefaultAtreusReqIDKey, "cvalue", "app", "test", "token", "test", "method", "normal", "key2", "value2"))
	cloneMd := extractIncomingAndClone(ctx)
	cloneMd.Set("key4", "value4")
	ctx = metadata.NewOutgoingContext(ctx, cloneMd)

	fromOut, _ := metadata.FromOutgoingContext(ctx)
	fmt.Println("metadata.FromOutgoingContext(ctx)=", fromOut)
	for i := 0; i < 1; i++ {
		msg := fmt.Sprintf("hello golang:%d", i)
		resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: msg})
		if err != nil {
			fmt.Println(err)
			fail++
		} else {
			fmt.Printf("normal hello: resp=%v, error=%v\n", resp, err)
		}
	}
	fmt.Println(fail)
}
