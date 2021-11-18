package retrys

import (
	"context"
	"time"

	"grpc-wrapper-framework/logger"

	"go.uber.org/zap"
)

// attempt: 当前重试次数
func WaitRetryBackoff(ctx context.Context, attempt int, callOpts *retryOptions) error {
	var (
		waitTime time.Duration = 0
	)
	if attempt > 0 {
		//获取重试等待时间
		waitTime = callOpts.backoffFunc(ctx, attempt)
	}

	loger := logger.FromContext(ctx)

	if waitTime > 0 {
		if loger != nil {
			loger.Info("WaitRetryBackoff caller for sleeping...", zap.Int("attempt", attempt), zap.Any("backoff duration", waitTime))
		}

		//set a block timer
		timer := time.NewTimer(waitTime)
		select {
		case <-ctx.Done():
			timer.Stop()
			return contextErr2GrpcErr(ctx.Err())
		case <-timer.C:
			//wake up
		}
	}
	return nil
}
