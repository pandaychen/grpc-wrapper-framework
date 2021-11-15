package config

import "time"

type BreakerConfig struct {
	On                         bool          `json:"on-off"`
	BreakerType                string        `json:"type"` //熔断器类型
	MaxRequestsForHalfOpen     int           `json:"max_request"`
	Interval                   time.Duration `json:"interval"`
	TimeoutForOpen             time.Duration `json:"timeout"`
	ReadyToTripForTotalrequets int           `json:"r2t_total_request"` //ReadyToTrip中统计的总请求次数
	ReadyToTripForFailratio    float64       `json:"r2t_fail_ratio"`    //ReadyToTrip中统计的出错比率
}

/*
//https://github.com/sony/gobreaker
type Settings struct {
	Name          string
	MaxRequests   uint32        // 半开状态期最大允许放行请求：即进入Half-Open状态时，一个时间周期内允许最大同时请求数（如果还达不到切回closed状态条件，则不能再放行请求）
	Interval      time.Duration // closed状态时，重置计数的时间周期；如果配为0，切入Open后永不切回Closed（不建议设置此值）
	Timeout       time.Duration // 进入Open状态后，多长时间会自动切成 Half-open，默认60s，不能配为0

    // ReadyToTrip回调函数：进入Open状态的条件，比如默认是连接5次出错，即进入Open状态，即可对熔断条件进行配置。在fail计数发生后，回调一次
	ReadyToTrip   func(counts Counts) bool

	// 状态切换时的熔断器
	OnStateChange func(name string, from State, to State)
}
*/
