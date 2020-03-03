package atreus

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//计时(测试interceptor顺序)
func Middle(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	fmt.Println("Call FakeRpcLog interceptor..")

	return handler(ctx, req)
}
