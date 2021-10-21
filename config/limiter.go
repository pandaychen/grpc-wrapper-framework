package config

type LimiterConfig struct {
	LimiterType string `json:"limiter_type"`
	LimiterRate int    `json:"limiter_rate"`
	LimiterSize int    `json:"limiter_size"`
}
