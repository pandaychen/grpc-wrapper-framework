package config

type LimiterConfig struct {
	On          bool   `json:"on"`
	LimiterType string `json:"type"`
	LimiterRate int    `json:"rate"`
	LimiterSize int    `json:"bucketsize"`
}
