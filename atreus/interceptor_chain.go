package atreus

//实现Unary interceptor-chain
//参考：https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/chain.go

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func AtreusUnaryInterceptorChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		chain := handler // 原始的RPC方法
		for i := len(interceptors) - 1; i >= 0; i-- {
			//从数组最后一个interceptor开始，依次和前一个建立chain关系
			chain = create_subchain(interceptors[i], chain, info)
		}

		//返回的是第0号位置上的chain
		return chain(ctx, req)
	}
}

func create_subchain(us_interceptor grpc.UnaryServerInterceptor, us_handler grpc.UnaryHandler, info *grpc.UnaryServerInfo) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return us_interceptor(ctx, req, info, us_handler)
	}
}
