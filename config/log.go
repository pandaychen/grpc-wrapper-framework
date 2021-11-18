package config

type LogConfig struct {
	//Level      string `json:"level"`
	FileName   string  `json:"file_name"`
	MaxSize    int     `json:"max_size"`
	MaxBackups int     `json:"max_backups"`
	MaxAge     int     `json:"max_age"`
	Compress   bool    `json:"compress"`
	Sampling   float64 `json:"sampling"`
}
