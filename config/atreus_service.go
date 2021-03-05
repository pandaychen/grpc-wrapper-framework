package config

import (
	"errors"
	"time"
)

type AtreusSvcConfig struct {
	Addr              string        `json:"address"`
	Keepalive         bool          `json:"keepalive"`
	Timeout           time.Duration `json:"timeout"`
	IdleTimeout       time.Duration `json:"idle_timeout"`
	MaxLifeTime       time.Duration `json:"max_life"`
	ForceCloseWait    time.Duration `json:"close_wait"`
	KeepAliveInterval time.Duration `json:"keepalive_interval"`
	KeepAliveTimeout  time.Duration `json:"keepalive_timeout"`

	//TLS config
	TLSon     bool   `json:"tls_on"`
	TLSCert   string `json:"tls_cert"`
	TLSKey    string `json:"tls_key"`
	TLSCaCert string `json:"tls_ca_cert"`

	//Service Register
	RegisterType        string        `json:"reg_type"`
	RegisterEndpoints   string        `json:"reg_endpoint"`
	RegisterTTL         time.Duration `json:"reg_ttl"`
	RegisterAPIOn       bool          `json:"reg_api_on"`
	RegisterRootPath    string        `json:"reg_root_path"`
	RegisterService     string        `json:"reg_service_name"`
	RegisterServiceVer  string        `json:"reg_service_version"`
	RegisterServiceAddr string        `json:"reg_service_addr"`

	//Limiter
	LimiterType string `json:"limiter_type"`
	LimiterRate int    `json:"limiter_rate"`
	LimiterSize int    `json:"limiter_size"`

	EtcdConfig
}

//global
var atreus_svc_config AtreusSvcConfig

func GetAtreusSvcConfig() *AtreusSvcConfig {
	return &atreus_svc_config
}

func AtreusSvcConfigInit() {
	Config := vipers.Use("atreus")
	if Config == nil {
		panic(errors.New("find grpc service config error"))
		return
	}
	SubconfigServer := Config.Use("server")
	if SubconfigServer == nil {
		panic(errors.New("find grpc service config error"))
		return
	}

	atreus_svc_config.Addr = SubconfigServer.GetString("address")
	atreus_svc_config.Keepalive = SubconfigServer.MustBool("keepalive", false)
	atreus_svc_config.Timeout = SubconfigServer.MustDuration("timeout", time.Second*10)
	atreus_svc_config.IdleTimeout = SubconfigServer.MustDuration("idle_timeout", time.Second*10)
	atreus_svc_config.MaxLifeTime = SubconfigServer.MustDuration("max_life", time.Second*10)
	atreus_svc_config.ForceCloseWait = SubconfigServer.MustDuration("close_wait", time.Second*10)
	atreus_svc_config.KeepAliveInterval = SubconfigServer.MustDuration("keepalive_interval", time.Second*10)
	atreus_svc_config.KeepAliveTimeout = SubconfigServer.MustDuration("keepalive_timeout", time.Second*10)
	atreus_svc_config.TLSon = SubconfigServer.MustBool("tls_on", false)
	atreus_svc_config.TLSCert = SubconfigServer.GetString("tls_cert")
	atreus_svc_config.TLSKey = SubconfigServer.GetString("tls_key")
	atreus_svc_config.TLSCaCert = SubconfigServer.GetString("tls_ca_cert")
	atreus_svc_config.RegisterType = SubconfigServer.MustString("reg_type", "etcd")
	atreus_svc_config.RegisterEndpoints = SubconfigServer.MustString("reg_endpoint", "http://127.0.0.1:2379")
	atreus_svc_config.RegisterTTL = SubconfigServer.MustDuration("reg_ttl", 10*time.Second)
	atreus_svc_config.RegisterAPIOn = SubconfigServer.MustBool("reg_api_on", false)
	atreus_svc_config.RegisterRootPath = SubconfigServer.MustString("reg_root_path", "/")
	atreus_svc_config.RegisterService = SubconfigServer.MustString("reg_service_name", "test")
	atreus_svc_config.RegisterServiceVer = SubconfigServer.MustString("reg_service_version", "v1.0")
	atreus_svc_config.RegisterServiceAddr = SubconfigServer.MustString("reg_service_addr", atreus_svc_config.Addr)

	atreus_svc_config.LimiterType = SubconfigServer.GetString("limiter_type")
	atreus_svc_config.LimiterRate = SubconfigServer.GetInt("limiter_rate")
	atreus_svc_config.LimiterSize = SubconfigServer.GetInt("limiter_size")
}

/*
func main() {
	InitConfigAbpath("./", "grpc_server", "yaml")
	AtreusSvcConfigInit()
	fmt.Println(GetAtreusSvcConfig())
}
*/
