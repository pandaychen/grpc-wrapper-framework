package atreus

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const slowThreshold = time.Millisecond * 500

//计时(最后一个拦截器)
func TimingOld(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	fmt.Println("Call RpcLog interceptor..")
	fmt.Printf("rpc=%s, req=%v", info.FullMethod, req)

	//final call rpc（if there is no interceptor） and get result
	resp, err = handler(ctx, req)
	fmt.Printf("finished %s, took=%v, resp=%v, err=%v", info.FullMethod, time.Since(start), resp, err)

	return resp, err
}

// Timing is an interceptor that logs the processing time (for client)
func (c *Client) Timing() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		serverName := path.Join(cc.Target(), method)

		//开始时间
		startTime := time.Now()

		err = invoker(ctx, method, req, reply, cc, opts...)
		elapseTime := time.Since(startTime)
		if err != nil {
			c.Logger.Error("[Client]Timing RPC Call fail", zap.Any("elapse time", elapseTime), zap.String("servername", serverName), zap.Any("req", req), zap.String("errmsg", err.Error()))
			return err
		}

		if elapseTime > slowThreshold {
			c.Logger.Info("[Client]Timing RPC Call Slow", zap.Any("elapse time", elapseTime), zap.String("servername", serverName), zap.Any("req", req))
		}
		return
	}
}

// 服务端接口调用耗时拦截器
func (s *Server) Timing() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)
			s.logTiming(ctx, args.FullMethod, req, duration)
		}()

		resp, err = handler(ctx, req)
		metricRpcServerReqDuration.Observe(float64(time.Since(startTime)/time.Millisecond), args.FullMethod)
		return
	}
}

func (s *Server) logTiming(ctx context.Context, method string, req interface{}, elapse_time time.Duration) {
	var addr string
	client, ok := peer.FromContext(ctx)
	if ok {
		//extractor src addr
		addr = client.Addr.String()
	}
	retstr, err := json.Marshal(req)
	if err != nil {
		s.Logger.Error("[Server]Timing RPC Call fail", zap.Any("elapse time", elapse_time), zap.String("servername", method), zap.String("errmsg", err.Error()), zap.String("callip", addr))
	}
	if elapse_time > slowThreshold {
		s.Logger.Info("[Server]Timing RPC Call Slow", zap.Any("elapse time", elapse_time), zap.String("servername", method), zap.String("callip", addr))
	}
	s.Logger.Info("[Server]Timing RPC Call", zap.Any("elapse time", elapse_time), zap.String("servername", method), zap.String("callip", addr), zap.String("content", string(retstr)))
}
