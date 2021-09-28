package codes

import (
	"sync"

	"github.com/golang/protobuf/ptypes/any"
	pbstatus "google.golang.org/genproto/googleapis/rpc/status"
	rpcodes "google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

var (
	GlobalAtreusCodesStore sync.Map //全局错误码数组
	//
	PB_SUCC = &pbstatus.Status{
		Code:    int32(rpcodes.OK),
		Message: "SUCC",
		Details: make([]*any.Any, 0),
	}

	PB_STAT_SUCC = &PbStatusWrapper{
		Status: PB_SUCC,
	}
)

// Breaker：Acceptable checks if given error is acceptable
func CheckErrorIsAcceptable(err error) bool {
	switch status.Code(err) {
	case rpcodes.DeadlineExceeded, rpcodes.Internal, rpcodes.Unavailable, rpcodes.DataLoss:
		//gRPC内部错误，过滤掉
		return false
	default:
		//逻辑错误
		return true
	}
}
