package atreus

import (
	"github.com/pandaychen/grpc-wrapper-framework/config"
	"github.com/pandaychen/grpc-wrapper-framework/pkg/k8s"
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

	//Register
	RegisterType      string         `json:"reg_type"`
	RegisterEndpoints string         `json:"reg_endpoint"`
	RegisterTTL       xtime.Duration `json:"reg_ttl"`
	RegisterAPIOn     bool           `json:"reg_api_on"`
}

func NewAtreusServerConfig() *AtreusServerConfig {
	//return default config
	return &AtreusServerConfig{}
}

//validate and generate svc config
func NewAtreusServerConfig2(conf *config.AtreusSvcConfig) *AtreusServerConfig {
	if config == nil {
		return
	} else {
		config := &AtreusServerConfig{
			Addr:              conf.Addr,
			Timeout:           xtime.Duration(config.Timeout),
			IdleTimeout:       xtime.Duration(conf.IdleTimeout),
			MaxLifeTime:       conf.MaxLifeTime,
			ForceCloseWait:    conf.ForceCloseWait,
			KeepAliveInterval: conf.KeepAliveInterval,
			KeepAliveTimeout:  conf.KeepAliveTimeout,
			TLSon:             conf.TLSon,
			TLSCert:           conf.TLSCert,
			TLSKey:            conf.TLSKey,
			TLSCaCert:         conf.TLSCaCert,
			RegisterType:      conf.RegisterType,
			RegisterEndpoints: conf.RegisterEndpoints,
			RegisterTTL:       conf.RegisterTTL,
			RegisterAPIOn:     conf.RegisterAPIOn,
		}

		return config
	}
}

// 构建K8S-ENV配置
func InitAtreusServerConfigK8S() (*AtreusServerConfig, error) {
	config := new(AtreusServerConfig)
	k8s.IgnorePrefix()
	err := k8s.FillConfig(config)
	if err != nil {
		panic(err)
		return nil, err
	}
	return config, nil
}
