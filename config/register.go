package config

import "time"

type RegistryConfig struct {
	RegisterType        string        `json:"reg_type"`
	RegisterEndpoints   string        `json:"reg_endpoint"`
	RegisterTTL         time.Duration `json:"reg_ttl"`
	RegisterAPIOn       bool          `json:"reg_api_on"`
	RegisterRootPath    string        `json:"reg_root_path"`
	RegisterService     string        `json:"reg_service_name"`
	RegisterServiceVer  string        `json:"reg_service_version"`
	RegisterServiceAddr string        `json:"reg_service_addr"`
}
