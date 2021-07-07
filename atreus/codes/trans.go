package codes

import (
	"github.com/golang/protobuf/ptypes/any"
	pbstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
)

//将err全部转换为grpc的错误码处理，如果错误码不属于grpc定义，则返回codes.Unknown/原始错误信息，否则直接返回grpc内置codes封装的wrapper结构

//https://github.com/grpc/grpc-go/blob/v1.36.1/status/status.go#L81
func TransformError2Codes(err error) *PbStatusWrapper {
	if err == nil {
		return PB_STAT_SUCC
	}

	rpc_error, _ := status.FromError(err)
	return &PbStatusWrapper{
		&pbstatus.Status{
			Code:    int32(rpc_error.Code()),
			Message: rpc_error.Message(),
			Details: make([]*any.Any, 0),
		},
	}
}
