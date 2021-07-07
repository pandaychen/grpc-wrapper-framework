package atreus

import (
	"github.com/pandaychen/grpc-wrapper-framework/config"
	"github.com/pandaychen/grpc-wrapper-framework/pkg/k8s"
	"github.com/pandaychen/grpc-wrapper-framework/pkg/xtime"
)

type AtreusServerConfig struct {
	Addr              string         `json:"address"`
	Keepalive         bool           `json:"keepalive"`
	Timeout           xtime.Duration `json:"timeout"`
	IdleTimeout       xtime.Duration `json:"idle_timeout"`
	MaxLifeTime       xtime.Duration `json:"max_life"`
	ForceCloseWait    xtime.Duration `json:"close_wait"`
	KeepAliveInterval xtime.Duration `json:"keepalive_interval"`
	KeepAliveTimeout  xtime.Duration `json:"keepalive_timeout"`
	InitWeight        string         `json:"init_weight"`

	//TLS config
	TLSon     bool   `json:"tls_on"`
	TLSCert   string `json:"tls_cert"`
	TLSKey    string `json:"tls_key"`
	TLSCaCert string `json:"tls_ca_cert"`

	//Register
	Regon               bool           `json:"reg_on"`
	RegisterType        string         `json:"reg_type"`
	RegisterEndpoints   string         `json:"reg_endpoint"`
	RegisterTTL         xtime.Duration `json:"reg_ttl"`
	RegisterAPIOn       bool           `json:"reg_api_on"`
	RegisterRootPath    string         `json:"reg_root_path"`
	RegisterService     string         `json:"reg_service_name"`
	RegisterServiceVer  string         `json:"reg_service_version"`
	RegisterServiceAddr string         `json:"reg_service_addr"`
}

func NewAtreusServerConfig() *AtreusServerConfig {
	//return default config
	return &AtreusServerConfig{}
}

//validate and generate svc config
func NewAtreusServerConfig2(conf *config.AtreusSvcConfig) *AtreusServerConfig {
	if conf == nil {
		return nil
	} else {
		config := &AtreusServerConfig{
			Addr:                conf.Addr,
			Keepalive:           conf.Keepalive,
			Timeout:             xtime.Duration(conf.Timeout),
			IdleTimeout:         xtime.Duration(conf.IdleTimeout),
			MaxLifeTime:         xtime.Duration(conf.MaxLifeTime),
			ForceCloseWait:      xtime.Duration(conf.ForceCloseWait),
			KeepAliveInterval:   xtime.Duration(conf.KeepAliveInterval),
			KeepAliveTimeout:    xtime.Duration(conf.KeepAliveTimeout),
			TLSon:               conf.TLSon,
			TLSCert:             conf.TLSCert,
			TLSKey:              conf.TLSKey,
			TLSCaCert:           conf.TLSCaCert,
			RegisterType:        conf.RegisterType,
			RegisterEndpoints:   conf.RegisterEndpoints,
			RegisterTTL:         xtime.Duration(conf.RegisterTTL),
			RegisterAPIOn:       conf.RegisterAPIOn,
			RegisterService:     conf.RegisterService,
			RegisterServiceVer:  conf.RegisterServiceVer,
			RegisterServiceAddr: conf.RegisterServiceAddr,
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

//RPC客户端配置封装
type AtreusClientConfig struct {
	DialTimeout  xtime.Duration `json:"dial_timeout"`
	Timeout      xtime.Duration `json:"timeout"`
	NonBlockSign bool           `json:"non_block_sign"` //是否默认阻塞

	//客户端服务发现
	k8sSign     bool   `json:"k8s_sign"`
	DnsSign     bool   `json:"dns_sign"`
	ServiceName string `json:"service_name"`
	ServicePort int    `json:"service_port"`
	DialScheme  string `json:"dial_scheme"`

	//注册中心
	RegisterType       string `json:"reg_type"`
	RegisterEndpoints  string `json:"reg_endpoint"`
	RegisterRootPath   string `json:"reg_root_path"`
	RegisterService    string `json:"reg_service_name"`
	RegisterServiceVer string `json:"reg_service_version"`

	//TLS config
	TLSon     bool   `json:"tls_on"`
	TLSCert   string `json:"tls_cert"`
	TLSKey    string `json:"tls_key"`
	TLSCaCert string `json:"tls_ca_cert"`

	//keepalive配置
	KeepaliveInterval      xtime.Duration `json:"keepalive_interval"`
	KeepaliveTimeout       xtime.Duration `json:"keepalive_timeout"`
	KeepaliveWithoutStream bool           `json:"keepalive_without_stream"`
}

//validate and generate svc config
func NewAtreusClientConfig2(conf *config.AtreusSvcConfig) *AtreusClientConfig {
	if conf == nil {
		return nil
	} else {
		config := &AtreusClientConfig{
			NonBlockSign:           conf.Keepalive,
			DialTimeout:            xtime.Duration(conf.Timeout),
			Timeout:                xtime.Duration(conf.IdleTimeout),
			KeepaliveInterval:      xtime.Duration(conf.KeepAliveInterval),
			KeepaliveTimeout:       xtime.Duration(conf.KeepAliveTimeout),
			KeepaliveWithoutStream: false,
			//客户端服务发现
			k8sSign:            false,
			DnsSign:            false,
			ServiceName:        "test-service",
			ServicePort:        8080,
			DialScheme:         "etcdv3",
			TLSon:              conf.TLSon,
			TLSCert:            conf.TLSCert,
			TLSKey:             conf.TLSKey,
			TLSCaCert:          conf.TLSCaCert,
			RegisterType:       conf.RegisterType,
			RegisterEndpoints:  conf.RegisterEndpoints,
			RegisterService:    conf.RegisterService,
			RegisterServiceVer: conf.RegisterServiceVer,
		}

		return config
	}
}
