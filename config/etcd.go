package config

import "time"

type EtcdConfig struct {
	Endpoints      []string      `json:"endpoints"`
	ConnectTimeout time.Duration `json:"timeout"`
	Secure         bool          `json:"secure"`
	TTL            int           `json:"ttl"`

	//Etcd
	DialKeepAliveTime    time.Duration `json:"dialkeepalivetime"`
	DialKeepAliveTimeout time.Duration `json:"dialkeepalivetimeout"`

	//ETCD 认证参数
	CertFilePath string `json:"certfilepath"`
	KeyFilePath  string `json:"keyfilepath"`
	CaCertPath   string `json:"cacertpath"`
	BasicAuth    bool   `json:"basicauth"`
	UserName     string `json:"username"`
	Password     string `json:"passwd"`
}
