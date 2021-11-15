package config

type TlsConfig struct {
	TLSon         bool   `json:"on-off"`
	TlsCommonName string `json:"cert_name"` //服务端证书name
	TLSCert       string `json:"tls_cert"`
	TLSKey        string `json:"tls_key"`
	TLSCaCert     string `json:"tls_ca_cert"`
}
