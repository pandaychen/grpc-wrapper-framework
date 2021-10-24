package atreus

import (
	"context"
	xmd "grpc-wrapper-framework/microservice/metadata"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//将cpu数据作为拦截器，每一次rpc调用都采集并返回客户端
/*
	INFO 05/18-06:44:36.358 grpc-access-log ret=0 path=/testproto.Greeter/SayHello ts=0.000648521 args=name:"tom" age:23  ip=127.0.0.1:8081
	get reply: {hello tom from 127.0.0.1:8081 false} map[cpu_usage:[36] serverinfo:[enjoy]]
*/
func (s *Server) ServerStat() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var (
			cpu_usage int64
		)
		resp, err = handler(ctx, req)
		//在rpc完成时采集，不区分错误
		//请求处理完成时，设置cpu-stat，读取一次cpu的瞬时值
		if cpu_usage != 0 {
			//每次客户端RPC请求，服务端都会计算cpu使用率，gRPC客户端在Pick的DoneInfo中获取此值
			trailer := metadata.Pairs([]string{xmd.CPUloadKey, strconv.FormatInt(int64(cpu_usage), 10)}...)
			//每次rpc请求时，放在tailer，返回给客户端
			grpc.SetTrailer(ctx, trailer)
		}
		return
	}
}
