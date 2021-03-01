package xtime

//from KRATOS，主要用于mysql时间戳转换/配置文件读取并转换/Context超时时间比较
import (
	xtime "time"
)

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration xtime.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := xtime.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}
