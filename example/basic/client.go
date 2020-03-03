package main

import (
	"fmt"

	pb "github.com/pandaychen/grpc-wrapper-framework/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:50001", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterServiceClient(conn)
	// ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("ckey", "cvalue"))
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "hello golang"})
	fmt.Printf("normal hello: resp=%v, error=%v\n", resp, err)
}
