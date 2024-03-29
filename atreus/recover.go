package atreus

import (
	"fmt"
	"grpc-wrapper-framework/errcode"
	"os"
	"runtime"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MAX_STACK_SIZE = 2048
)

// Recovery interceptor：必须放在第 0 号链位置
func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//user defer to recovery
	defer func() {
		// 定义堆栈恢复逻辑（打印 coredump 时的堆栈信息）
		if r := recover(); r != nil {
			stack := make([]byte, MAX_STACK_SIZE)
			stack = stack[:runtime.Stack(stack, false)]
			fmt.Printf("Panic Rpc Call: %s, err=%v, stack:\n%s", info.FullMethod, r, string(stack))
			err = status.Errorf(codes.Unknown, "Server panic error: %s", info.FullMethod)
		}
	}()

	// 这里返回的是下一个（interceptor）链
	return handler(ctx, req)
}

// NICE：将 recovery 作为 Server 拦截器，调用，打印崩溃异常的 stack 信息
// 在 Server 初始化时，这样调用，s.Use(s.Recovery(),...)
func (s *Server) Recovery() grpc.UnaryServerInterceptor {
	// 服务端的 rpc == handler
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				rs := runtime.Stack(buf, false)
				if rs > size {
					rs = size
				}
				buf = buf[:rs]
				pl := fmt.Sprintf("Panic Rpc Call: : %v\n%v\n%s\n", req, rerr, buf)
				fmt.Fprintf(os.Stderr, pl)
				//err = status.Errorf(codes.Unknown, "Server panic error: %s", args.FullMethod)

				//返回grpc的内置错误码+项目的错误码封装
				err = status.Errorf(codes.Unknown, errcode.ServerErr.Error())
			}
		}()
		// 注意：服务端的拦截器 handler，这里是进入下一个拦截器
		resp, err = handler(ctx, req)
		return
	}
}

// 客户端 recovery 拦截器
func (c *Client) Recovery() grpc.UnaryClientInterceptor {
	// 客户端 == invoker
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				rs := runtime.Stack(buf, false)
				if rs > size {
					rs = size
				}
				buf = buf[:rs]
				pl := fmt.Sprintf("client panic: %v\n%v\n%v\n%s\n", req, reply, rerr, buf)
				fmt.Fprintf(os.Stderr, pl)
				err = errcode.ServerErr
			}
		}()
		// 注意：客户端的拦截器 invoker，这里是进入下一个拦截器
		err = invoker(ctx, method, req, reply, cc, opts...)
		return
	}
}
