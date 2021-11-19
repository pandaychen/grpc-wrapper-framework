package retrys

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

//https://github.com/kasvith/simplelb/blob/master/main.go
//https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/retry

// retry 配置
type retryOptions struct {
	maxRetries         int           //
	perCallTimeout     time.Duration //WARNING! 对每次重试请求的单独设置超时时间
	includeRetryHeader bool
	codes              []codes.Code
	backoffFunc        BackoffFuncContext
}

func (o *retryOptions) GetMaxretries() int {
	return o.maxRetries
}

func (o *retryOptions) GetperCallTimeout() time.Duration {
	return o.perCallTimeout
}

// CallOption is a grpc.CallOption that is local to grpc_retry.
type CallOption struct {
	grpc.EmptyCallOption // make sure we implement private after() and before() fields so we don't panic.

	//设置 retryOptions的方法
	applyFunc func(opt *retryOptions)
}

// 按照callOptions的配置更新opt的原始配置
func ReuseOrNewWithCallOptions(opt *retryOptions, callOptions []CallOption) *retryOptions {
	if len(callOptions) == 0 {
		//if callOptions not set,return opt
		return opt
	}
	optCopy := &retryOptions{}
	*optCopy = *opt
	for _, f := range callOptions {
		f.applyFunc(optCopy)
	}
	return optCopy
}

//从原始的grpc.callOptions数组配置中剥离grpcOptions与retryOptions
func SplitCallOptions(callOptions []grpc.CallOption) (grpcOptions []grpc.CallOption, retryOptions []CallOption) {
	for _, opt := range callOptions {
		if co, ok := opt.(CallOption); ok {
			retryOptions = append(retryOptions, co)
		} else {
			grpcOptions = append(grpcOptions, opt)
		}
	}
	return grpcOptions, retryOptions
}

/////////////retryOptions设置方法

// 设置最大的重传次数
func WithMax(maxRetries int) CallOption {
	return CallOption{
		applyFunc: func(o *retryOptions) {
			o.maxRetries = maxRetries
		}}
}

// 设置每次重传时的backoff时间计算回调方法，系统内置的方法见grpc-wrapper-framework/microservice/retrys/backoff.go
func WithBackoff(bkf BackoffFunc) CallOption {
	return CallOption{applyFunc: func(o *retryOptions) {
		o.backoffFunc = BackoffFuncContext(func(ctx context.Context, attempt int) time.Duration {
			return bkf(attempt)
		})
	}}
}

//
func WithBackoffContext(bkfc BackoffFuncContext) CallOption {
	return CallOption{applyFunc: func(o *retryOptions) {
		o.backoffFunc = bkfc
	}}
}

// WithCodes sets which codes should be retried.
//
// Please *use with care*, as you may be retrying non-idempotent calls.
//
// You cannot automatically retry on Cancelled and Deadline, please use `WithPerRetryTimeout` for these.
func WithCodes(retryCodes ...codes.Code) CallOption {
	return CallOption{applyFunc: func(o *retryOptions) {
		o.codes = retryCodes
	}}
}

func WithHeaderSignOff(on bool) CallOption {
	return CallOption{applyFunc: func(o *retryOptions) {
		o.includeRetryHeader = on
	}}
}

// WithPerRetryTimeout sets the RPC timeout per call (including initial call) on this call, or this interceptor.
//
// The context.Deadline of the call takes precedence and sets the maximum time the whole invocation
// will take, but WithPerRetryTimeout can be used to limit the RPC time per each call.
//
// For example, with context.Deadline = now + 10s, and WithPerRetryTimeout(3 * time.Seconds), each
// of the retry calls (including the initial one) will have a deadline of now + 3s.
//
// A value of 0 disables the timeout overrides completely and returns to each retry call using the
// parent `context.Deadline`.
//
// Note that when this is enabled, any DeadlineExceeded errors that are propagated up will be retried.
func WithPerRetryTimeout(timeout time.Duration) CallOption {
	return CallOption{applyFunc: func(o *retryOptions) {
		o.perCallTimeout = timeout
	}}
}

// retryOptions默认配置值

var (
	// DefaultRetriableCodes is a set of well known types gRPC codes that should be retri-able.
	//
	// `ResourceExhausted` means that the user quota, e.g. per-RPC limits, have been reached.
	// `Unavailable` means that system is currently unavailable and the client should retry again.
	DefaultRetriableCodes = []codes.Code{codes.ResourceExhausted, codes.Unavailable}

	DefaultOptions = &retryOptions{
		maxRetries:         0,                // disabled
		perCallTimeout:     time.Duration(0), // disabled
		includeRetryHeader: true,
		codes:              DefaultRetriableCodes,
		backoffFunc: BackoffFuncContext(func(ctx context.Context, attempt int) time.Duration {
			return BackoffWithLinearJitter(50*time.Millisecond /*jitter*/, 0.1)(attempt)
		}),
	}
)
