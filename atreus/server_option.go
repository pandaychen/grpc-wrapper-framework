package atreus

import (
	"github.com/pandaychen/grpc-wrapper-framework/config"
	"github.com/pandaychen/grpc-wrapper-framework/pkg/xtime"
)

type AtreusServerConfig struct {
	Addr              string         `json:"address"`
	Timeout           xtime.Duration `json:"timeout"`
	IdleTimeout       xtime.Duration `json:"idle_timeout"`
	MaxLifeTime       xtime.Duration `json:"max_life"`
	ForceCloseWait    xtime.Duration `json:"close_wait"`
	KeepAliveInterval xtime.Duration `json:"keepalive_interval"`
	KeepAliveTimeout  xtime.Duration `json:"keepalive_timeout"`

	//TLS config
	TLSon     bool   `json:"tls_on"`
	TLSCert   string `json:"tls_cert"`
	TLSKey    string `json:"tls_key"`
	TLSCaCert string `json:"tls_ca_cert"`

	//register
	RegisterType      string         `json:"reg_type"`
	RegisterEndpoints string         `json:"reg_endpoint"`
	RegisterTTL       xtime.Duration `json:"reg_ttl"`
	RegisterAPIOn     bool           `json:"reg_api_on"`
}

func NewAtreusServerConfig() *AtreusServerConfig {
	//return default config
	return &AtreusServerConfig{}
}

func NewAtreusServerConfig2(conf *config.AtreusSvcConfig) *AtreusServerConfig {
	/*
		config := &AtreusServerConfig{
			Addr: Addr
		Timeout           xtime.Duration `json:"timeout"`
		IdleTimeout       xtime.Duration `json:"idle_timeout"`
		MaxLifeTime       xtime.Duration `json:"max_life"`
		ForceCloseWait    xtime.Duration `json:"close_wait"`
		KeepAliveInterval xtime.Duration `json:"keepalive_interval"`
		KeepAliveTimeout  xtime.Duration `json:"keepalive_timeout"`

		//TLS config
		TLSon     bool   `json:"tls_on"`
		TLSCert   string `json:"tls_cert"`
		TLSKey    string `json:"tls_key"`
		TLSCaCert string `json:"tls_ca_cert"`

		//register
		RegisterType      string         `json:"reg_type"`
		RegisterEndpoints string         `json:"reg_endpoint"`
		RegisterTTL       xtime.Duration `json:"reg_ttl"`
		RegisterAPIOn     bool           `json:"reg_api_on"`
		}
	*/

	return config
}

// 构建K8S-ENV
func InitAtreusServerConfigK8S() (*AtreusServerConfig, error) {
	cfg := new(AtreusServerConfig)
	bf_k8s.IgnorePrefix()
	err := bf_k8s.FillConfig(cfg)
	if err != nil {
		panic(err)
		return nil, err
	}
	return cfg, nil
}
