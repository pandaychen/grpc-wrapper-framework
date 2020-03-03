package atreus

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//计时(最后一个拦截器)
func Timing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	fmt.Println("Call RpcLog interceptor..")
	fmt.Println("rpc=%s, req=%v", info.FullMethod, req)

	//final call rpc（if there is no interceptor） and get result
	resp, err = handler(ctx, req)
	fmt.Println("finished %s, took=%v, resp=%v, err=%v", info.FullMethod, time.Since(start), resp, err)

	return resp, err
}
