package pyerrors

import "errors"

var (
	ErrorBreakerOpenServiceUnavailable = errors.New("circuit breaker opening")
)

var (
	RatelimiterServiceReject = "ErrRatelimit"

	InternalError = "ErrInternal"

	AuthenticatorMissing = "ErrAuthenticatorMissing"

	TokenVerifyInvalid = "ErrTokenVerify"
)
