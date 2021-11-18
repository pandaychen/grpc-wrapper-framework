package retrys

import (
	"context"
	"time"

	bkf "grpc-wrapper-framework/microservice/backoff"
)

//https://github.com/kasvith/simplelb/blob/master/main.go
//https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/retry

// BackoffFunc denotes a family of functions that control the backoff duration between call retries.
//
// They are called with an identifier of the attempt, and should return a time the system client should
// hold off for. If the time returned is longer than the `context.Context.Deadline` of the request
// the deadline of the request takes precedence and the wait will be interrupted before proceeding
// with the next iteration.
type BackoffFunc func(attempt int) time.Duration

// BackoffFuncContext denotes a family of functions that control the backoff duration between call retries.
//
// They are called with an identifier of the attempt, and should return a time the system client should
// hold off for. If the time returned is longer than the `context.Context.Deadline` of the request
// the deadline of the request takes precedence and the wait will be interrupted before proceeding
// with the next iteration. The context can be used to extract request scoped metadata and context values.
type BackoffFuncContext func(ctx context.Context, attempt int) time.Duration

//---BackoffFunc的实例化方法---//

//1. 线性退避
func BackoffWithLinear(waitBetween time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		return waitBetween
	}
}

//2. 带Jitter的线性退避
// For example waitBetween=1s and jitter=0.10 can generate waits between 900ms and 1100ms.
func BackoffWithLinearJitter(waitBetween time.Duration, jitterFraction float64) BackoffFunc {
	return func(attempt int) time.Duration {
		return bkf.JitterAlgo(waitBetween, jitterFraction)
	}
}

//3. 指数退避
// BackoffExponential produces increasing intervals for each attempt.
// The scalar is multiplied times 2 raised to the current attempt. So the first
// retry with a scalar of 100ms is 100ms, while the 5th attempt would be 1.6s.
func BackoffWithExponential(scalar time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		return bkf.ExponentBase2Algo(scalar, attempt)
	}
}

//4. 指数退避+Jitter
// BackoffExponentialWithJitter creates an exponential backoff like BackoffExponential does, but adds jitter.
func BackoffExponentialWithJitter(scalar time.Duration, jitterFraction float64) BackoffFunc {
	return func(attempt int) time.Duration {
		real_scalar := bkf.ExponentBase2Algo(scalar, attempt)
		return bkf.JitterAlgo(real_scalar, jitterFraction)
	}
}
