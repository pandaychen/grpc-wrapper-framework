package config

import (
	"errors"
	"time"
)

type AtreusSvcConfig struct {
	Addr              string        `json:"address"`
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

	//register
	RegisterType      string        `json:"reg_type"`
	RegisterEndpoints string        `json:"reg_endpoint"`
	RegisterTTL       time.Duration `json:"reg_ttl"`
	RegisterAPIOn     bool          `json:"reg_api_on"`
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
	atreus_svc_config.Timeout = SubconfigServer.MustDuration("timeout", time.Second*10)
}

/*
func main() {
	InitConfigAbpath("./", "grpc_server", "yaml")
	AtreusSvcConfigInit()
	fmt.Println(GetAtreusSvcConfig())
}
*/
