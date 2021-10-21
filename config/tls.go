package config

type TlsConfig struct {
	TLSon     bool   `json:"tls_on"`
	TLSCert   string `json:"tls_cert"`
	TLSKey    string `json:"tls_key"`
	TLSCaCert string `json:"tls_ca_cert"`
}
