package backoff

//退避算法（重试），参考https://aws.amazon.com/cn/blogs/architecture/exponential-backoff-and-jitter/
//指数退避

import (
	"math/rand"
	"time"
)

type BackoffAlgo struct {
	name string
}

// JitterUp adds random jitter to the duration.
//
// This adds or subtracts time from the duration within a given jitter fraction.
// For example for 10s and jitter 0.1, it will return a time within [9s, 11s])
func JitterAlgo(duration time.Duration, jitter float64) time.Duration {
	multiplier := jitter * (rand.Float64()*2 - 1)
	return time.Duration(float64(duration) * (1 + multiplier))
}

// ExponentBase2 computes 2^(a-1) where a >= 1. If a is 0, the result is 0.
func ExponentBase2Algo(duration time.Duration, attempt int) time.Duration {
	return duration * time.Duration((1<<attempt)>>1)
}
