package retrys

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	actx "grpc-wrapper-framework/common/context"
	amd "grpc-wrapper-framework/microservice/metadata"
)

func contextErr2GrpcErr(err error) error {
	switch err {
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

// 判断是否为context错误，这个错误可能由上层触发或者deadline导致，不能视为retry重传的错误条件（场景：上层主动取消调用，那么重传也没有必要了）
func IsContextError(err error) bool {
	code := status.Code(err)
	return code == codes.DeadlineExceeded || code == codes.Canceled
}

func CheckIsRetriable(err error, callOpts *retryOptions) bool {
	errCode := status.Code(err)
	if IsContextError(err) {
		// context errors are not retriable based on user settings.
		return false
	}

	//默认的重传code：codes.ResourceExhausted, codes.Unavailable
	for _, code := range callOpts.codes {
		if code == errCode {
			return true
		}
	}
	return false
}

func PerCallContext(ctx context.Context, attempt int, callOpts *retryOptions) context.Context {
	if callOpts.perCallTimeout != 0 {
		//rpc超时时间
		ctx, _ = context.WithTimeout(ctx, callOpts.perCallTimeout)
	}

	//在header中添加CtxAttemptKey标识
	if attempt > 0 && callOpts.includeRetryHeader {
		//需要在header中设置attemp标志
		var (
			newmd metadata.MD
		)
		//mdClone := amd.ExtractOutgoing(ctx).Clone().Set(actx.CtxAttemptKey, fmt.Sprintf("%d", attempt))
		//ctx = mdClone.ToOutgoing(ctx)

		//get fully copy of origin client context
		newmd = amd.CloneClientOutgoingData(ctx)
		newmd.Set(actx.CtxAttemptKey, fmt.Sprintf("%d", attempt))
		ctx = metadata.NewOutgoingContext(ctx, newmd)
	}
	return ctx
}
