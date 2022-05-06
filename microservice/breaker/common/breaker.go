package common

// breaker的通用封装

// Acceptable is the func to check if the error can be accepted
type CheckIfAcceptableCaller func(err error) bool

// fallback
type FallBackCaller func(err error) error

// Breaker is a CircuitBreaker pattern.
type Breaker interface {
	Name() string
	State() int32                                                                                // 获取当前熔断器的状态
	Allow() error                                                                                //是否允许请求
	MarkSuccess()                                                                                //数值统计
	MarkFailed()                                                                                 //数值统计
	DoRequest(req func() error, fallback FallBackCaller, isAccept CheckIfAcceptableCaller) error //业务方法调用
}

// breaker status
const (
	// StateOpen when circuit breaker open, request not allowed, after sleep
	// some duration, allow one single request for testing the health, if ok
	// then state reset to closed, if not continue the step.
	StateOpen int32 = iota
	// StateClosed when circuit breaker closed, request allowed, the breaker
	// calc the succeed ratio, if request num greater request setting and
	// ratio lower than the setting ratio, then reset state to open.
	StateClosed
	// StateHalfopen when circuit breaker open, after slepp some duration, allow
	// one request, but not state closed.
	StateHalfopen
)
