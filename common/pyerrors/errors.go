package pyerrors

import "errors"

var (
	ErrorBreakerOpenServiceUnavailable = errors.New("circuit breaker opening")
)
