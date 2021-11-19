package config

import "time"

type ClientRetryConfig struct {
	On             bool          `json:"on-off"`
	Maxretry       int           `json:"max_retry"`          //客户端的最大重试次数
	PerCallTimeout time.Duration `json:"per_call_timeout"`   //每次RPC请求单独设置超时时间
	HeaderSign     bool          `json:"inject_header_sign"` //是否注入header字段
}
