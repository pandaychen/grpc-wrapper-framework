package config

type AclConfig struct {
	On          bool     `json:"on-off"`
	WhiteIpList []string `json:white_list`
}
