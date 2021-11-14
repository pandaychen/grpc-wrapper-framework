package atreus

import (
	"fmt"
	"net"
	"strings"

	nw "grpc-wrapper-framework/pkg/network"

	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

//获取调用端IP
func GetClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	} else if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}

	addSlice := strings.Split(pr.Addr.String(), ":")
	if addSlice[0] == "[" {
		return "localhost", nil
	}
	return addSlice[0], nil
}

func (s *Server) SrcIpFilter() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		srcip, err := GetClientIP(ctx)
		if err != nil {
			s.Logger.Error("SrcIpFilter GetClientIP error", zap.Any("errmsg", err))
			return nil, status.Error(codes.InvalidArgument, "srcip unknown")
		}

		if bret := nw.CheckIpCidr(srcip, s.CallerIp); bret == false {
			s.Logger.Error("[SrcIpFilter]src ip forbidden", zap.String("srcip", srcip))
			return nil, status.Error(codes.InvalidArgument, "srcip restrict")
		}

		//go to next interceptor
		return handler(ctx, req)
	}
}
