package atreus

import (
	"fmt"
	"runtime"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	MAX_STACK_SIZE = 2048
)

// Recovery interceptor：必须放在第0号链位置
func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//user defer to recovery
	fmt.Println("Call Recovery interceptor...")
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, MAX_STACK_SIZE)
			stack = stack[:runtime.Stack(stack, false)]
			fmt.Printf("Panin Rpc Call: %s, err=%v, stack:\n%s", info.FullMethod, r, string(stack))
		}
	}()

	//这里返回的是下一个（interceptor）链
	return handler(ctx, req)
}
