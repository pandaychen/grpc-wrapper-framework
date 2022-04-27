package config

import "time"

type SysCollectorConfig struct {
	Cputype           string        `json:"type"` //cvm or docker
	CollectorDuration time.Duration `json:"duration"`
	Multiply          int           `json:"multiply"` //cpu采集结果放大倍数
}
