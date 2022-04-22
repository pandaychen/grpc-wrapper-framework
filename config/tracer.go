package config

type TracingConfig struct {
	ServiceName string `json:"service_name"`
	Collector   string `json:"collector"` //收集地址
	TracerType  string `json:"type"`
}
