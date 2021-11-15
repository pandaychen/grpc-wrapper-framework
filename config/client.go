package config

import (
	"errors"
	"time"
)

type CliConfig struct {
	DialAddress string        `json:"dial_address"`
	DialPort    int           `json:"dial_port"`
	DialScheme  string        `json:"dial_scheme"`
	Env         string        `json:"env"`
	LbType      string        `json:lbtype`
	Timeout     time.Duration `json:"timeout"`
}

type AtreusCliConfig struct {
	//Cli config
	CliConf *CliConfig

	//Dial
	SrvDnsConf *SrvDnsConfig

	//TLS config
	TlsConf *TlsConfig

	//Service Register
	RegistryConf *RegistryConfig

	//Log
	LogConf *LogConfig

	//Auth
	AuthConf *AuthConfig

	//Breaker
	BreakerConf *BreakerConfig
}

//global
var atreus_cli_config AtreusCliConfig

func GetAtreusCliConfig() *AtreusCliConfig {
	//lock for hot reload?
	return &atreus_cli_config
}

func AtreusCliConfigInit() {
	Config := vipers.Use("atreus")
	if Config == nil {
		panic(errors.New("find grpc service client config error"))
		return
	}
	SubconfigClient := Config.Use("client")
	if SubconfigClient == nil {
		panic(errors.New("find grpc service client config error"))
		return
	}
	atreus_cli_config.CliConf = new(CliConfig)
	atreus_cli_config.CliConf.DialAddress = SubconfigClient.GetString("dial_address")
	atreus_cli_config.CliConf.DialPort = SubconfigClient.GetInt("dial_port")
	atreus_cli_config.CliConf.DialScheme = SubconfigClient.GetString("dial_scheme")
	atreus_cli_config.CliConf.Env = SubconfigClient.GetString("env")
	atreus_cli_config.CliConf.LbType = SubconfigClient.GetString("lbtype")
	atreus_cli_config.CliConf.Timeout = SubconfigClient.MustDuration("timeout", 10*time.Second)

	atreus_cli_config.SrvDnsConf = new(SrvDnsConfig)
	SubDnsconfig := Config.Use("dnsservice")
	if SubDnsconfig == nil {
		//not set
	} else {
		atreus_cli_config.SrvDnsConf.SrvName = SubDnsconfig.GetString("name")
		atreus_cli_config.SrvDnsConf.SrvPort = SubDnsconfig.GetInt("port")
	}

	atreus_cli_config.TlsConf = new(TlsConfig)
	SubTlsconfig := Config.Use("security")
	if SubTlsconfig == nil {
		//not set
	} else {
		atreus_cli_config.TlsConf.TLSon = SubTlsconfig.MustBool("on-off", false)
		atreus_cli_config.TlsConf.TlsCommonName = SubTlsconfig.GetString("cert_name")
		atreus_cli_config.TlsConf.TLSKey = SubTlsconfig.GetString("tls_key")
		atreus_cli_config.TlsConf.TLSCert = SubTlsconfig.GetString("tls_cert")
		atreus_cli_config.TlsConf.TLSCaCert = SubTlsconfig.GetString("tls_ca_cert")
	}

	atreus_cli_config.RegistryConf = new(RegistryConfig)
	SubRegconfig := Config.Use("discovery")
	if SubRegconfig == nil {
		//not set
	} else {
		atreus_cli_config.RegistryConf.RegOn = SubRegconfig.MustBool("reg_on", false)
		atreus_cli_config.RegistryConf.RegisterType = SubRegconfig.MustString("reg_type", "etcd")
		atreus_cli_config.RegistryConf.RegisterEndpoints = SubRegconfig.MustString("reg_endpoint", "http://127.0.0.1:2379")
		atreus_cli_config.RegistryConf.RegisterTTL = SubRegconfig.MustDuration("reg_ttl", 10*time.Second)
		atreus_cli_config.RegistryConf.RegisterAPIOn = SubRegconfig.MustBool("reg_api_on", false)
		atreus_cli_config.RegistryConf.RegisterRootPath = SubRegconfig.MustString("reg_root_path", "/")
		atreus_cli_config.RegistryConf.RegisterService = SubRegconfig.MustString("reg_service_name", "test")
		atreus_cli_config.RegistryConf.RegisterServiceVer = SubRegconfig.MustString("reg_service_version", "v1.0")
	}

	atreus_cli_config.LogConf = new(LogConfig)
	SubLogConfig := vipers.Use("log")
	if SubLogConfig == nil {
		//not set
	} else {
		atreus_cli_config.LogConf.FileName = SubLogConfig.MustString("file_name", "../log/server.log")
		atreus_cli_config.LogConf.MaxSize = SubLogConfig.MustInt("max_size", 100)
		atreus_cli_config.LogConf.MaxBackups = SubLogConfig.MustInt("max_backups", 10)
		atreus_cli_config.LogConf.MaxAge = SubLogConfig.MustInt("max_age", 20)
		atreus_cli_config.LogConf.Compress = SubLogConfig.MustBool("compress", true)
	}

	atreus_cli_config.AuthConf = new(AuthConfig)
	SubAuthconfig := Config.Use("auth")
	if SubAuthconfig == nil {
		//not set
	} else {
		atreus_cli_config.AuthConf.On = SubAuthconfig.MustBool("on-off", false)
	}

	atreus_cli_config.BreakerConf = new(BreakerConfig)
	SubBreakerconfig := Config.Use("breaker")
	if SubBreakerconfig == nil {
		//not set
	} else {
		atreus_cli_config.BreakerConf.On = SubBreakerconfig.MustBool("on-off", false)
	}
}
