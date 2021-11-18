package config

import (
	"errors"
	"fmt"
	"time"
)

type SrvConfig struct {
	Addr              string        `json:"address"`
	Keepalive         bool          `json:"keepalive"`
	Timeout           time.Duration `json:"timeout"`
	IdleTimeout       time.Duration `json:"idle_timeout"`
	MaxLifeTime       time.Duration `json:"max_life"`
	ForceCloseWait    time.Duration `json:"close_wait"`
	KeepAliveInterval time.Duration `json:"keepalive_interval"`
	KeepAliveTimeout  time.Duration `json:"keepalive_timeout"`
	MaxRetry          int           `json:"max_retry"`
}

type AtreusSvcConfig struct {
	//Server config
	SrvConf *SrvConfig `json:"srv_conf"`

	//TLS config
	TlsConf *TlsConfig `json:"tls_conf"`

	//Service Register
	RegistryConf *RegistryConfig `json:"registry_conf"`

	//Limiter
	LimiterConf *LimiterConfig `json:"limiter_conf"`
	//Etcd
	EtcdConf *EtcdConfig `json:"etcd_conf"`

	//Weight
	WeightConf *WeightConfig `json:"weight_conf"`

	//Log
	LogConf *LogConfig `json:"log_conf"`

	//Auth
	AuthConf *AuthConfig `json:"auth_conf"`

	//ACL
	AclConf *AclConfig `json:"acl_conf"`
}

//global
var atreus_svc_config AtreusSvcConfig

func GetAtreusSvcConfig() *AtreusSvcConfig {
	//lock for hot reload?
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
	atreus_svc_config.SrvConf = new(SrvConfig)
	atreus_svc_config.SrvConf.Addr = SubconfigServer.GetString("address")
	atreus_svc_config.SrvConf.Keepalive = SubconfigServer.MustBool("keepalive", false)
	atreus_svc_config.SrvConf.Timeout = SubconfigServer.MustDuration("timeout", time.Second*10)
	atreus_svc_config.SrvConf.IdleTimeout = SubconfigServer.MustDuration("idle_timeout", time.Second*10)
	atreus_svc_config.SrvConf.MaxLifeTime = SubconfigServer.MustDuration("max_life", time.Second*10)
	atreus_svc_config.SrvConf.ForceCloseWait = SubconfigServer.MustDuration("close_wait", time.Second*10)
	atreus_svc_config.SrvConf.KeepAliveInterval = SubconfigServer.MustDuration("keepalive_interval", time.Second*10)
	atreus_svc_config.SrvConf.KeepAliveTimeout = SubconfigServer.MustDuration("keepalive_timeout", time.Second*10)
	atreus_svc_config.SrvConf.MaxRetry = SubconfigServer.GetInt("max_retry")

	atreus_svc_config.TlsConf = new(TlsConfig)
	SubTlsconfig := Config.Use("security")
	if SubTlsconfig == nil {
		//not set
	} else {
		atreus_svc_config.TlsConf.TLSon = SubTlsconfig.MustBool("on-off", false)
		atreus_svc_config.TlsConf.TLSCert = SubTlsconfig.GetString("tls_cert")
		atreus_svc_config.TlsConf.TLSKey = SubTlsconfig.GetString("tls_key")
		atreus_svc_config.TlsConf.TLSCaCert = SubTlsconfig.GetString("tls_ca_cert")
	}

	atreus_svc_config.RegistryConf = new(RegistryConfig)
	SubRegconfig := Config.Use("register")
	if SubRegconfig == nil {
		//not set
	} else {
		atreus_svc_config.RegistryConf.RegOn = SubRegconfig.MustBool("on-off", true)
		atreus_svc_config.RegistryConf.RegisterType = SubRegconfig.MustString("reg_type", "etcd")
		atreus_svc_config.RegistryConf.RegisterEndpoints = SubRegconfig.MustString("reg_endpoint", "http://127.0.0.1:2379")
		atreus_svc_config.RegistryConf.RegisterTTL = SubRegconfig.MustDuration("reg_ttl", 10*time.Second)
		atreus_svc_config.RegistryConf.RegisterAPIOn = SubRegconfig.MustBool("reg_api_on", false)
		atreus_svc_config.RegistryConf.RegisterRootPath = SubRegconfig.MustString("reg_root_path", "/")
		atreus_svc_config.RegistryConf.RegisterService = SubRegconfig.MustString("reg_service_name", "test")
		atreus_svc_config.RegistryConf.RegisterServiceVer = SubRegconfig.MustString("reg_service_version", "v1.0")
		atreus_svc_config.RegistryConf.RegisterServiceAddr = SubRegconfig.MustString("reg_service_addr", atreus_svc_config.SrvConf.Addr)
	}

	atreus_svc_config.LimiterConf = new(LimiterConfig)
	SubLimiterconfig := Config.Use("limiter")
	if SubLimiterconfig == nil {
		//not set
	} else {
		atreus_svc_config.LimiterConf.On = SubLimiterconfig.MustBool("on-off", false)
		atreus_svc_config.LimiterConf.LimiterType = SubLimiterconfig.GetString("type")
		atreus_svc_config.LimiterConf.LimiterRate = SubLimiterconfig.GetInt("rate")
		atreus_svc_config.LimiterConf.LimiterSize = SubLimiterconfig.GetInt("bucketsize")
	}

	atreus_svc_config.WeightConf = new(WeightConfig)
	SubWeigthConfig := Config.Use("weight")
	if SubLimiterconfig == nil {
		//not set
	} else {
		atreus_svc_config.WeightConf.Weight = SubWeigthConfig.GetString("init")
	}

	atreus_svc_config.AuthConf = new(AuthConfig)
	SubAuthconfig := Config.Use("auth")
	if SubAuthconfig == nil {
		//not set
	} else {
		atreus_svc_config.AuthConf.On = SubAuthconfig.GetBool("on-off")
	}

	atreus_svc_config.AclConf = new(AclConfig)
	SubAclconfig := Config.Use("acl")
	if SubAclconfig == nil {
		//not set
	} else {
		atreus_svc_config.AclConf.On = SubAclconfig.GetBool("on-off")
		atreus_svc_config.AclConf.WhiteIpList = SubAclconfig.MustStringSlice("white_list", []string{"127.0.0.1/32"})
	}

	//fmt.Println(atreus_svc_config.AclConf.WhiteIpList)

	atreus_svc_config.LogConf = new(LogConfig)
	SubLogConfig := vipers.Use("log")
	if SubLogConfig == nil {
		//not set
	} else {
		atreus_svc_config.LogConf.FileName = SubLogConfig.MustString("file_name", "../log/server.log")
		atreus_svc_config.LogConf.MaxSize = SubLogConfig.MustInt("max_size", 100)
		atreus_svc_config.LogConf.MaxBackups = SubLogConfig.MustInt("max_backups", 10)
		atreus_svc_config.LogConf.MaxAge = SubLogConfig.MustInt("max_age", 20)
		atreus_svc_config.LogConf.Compress = SubLogConfig.MustBool("compress", true)
	}

}

func main() {
	InitConfigAbsolutePath("./", "server", "yaml")
	AtreusSvcConfigInit()
	fmt.Println(GetAtreusSvcConfig())
}
