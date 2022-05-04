package window

import (
	"time"

	"go.uber.org/zap"
)

type SliderWindowOption func(sw *SliderWindow)

func IgnoreCurrentBucket() SliderWindowOption {
	return func(w *SliderWindow) {
		w.ignoreCurrentBucket = true
	}
}

// 设置滑动窗口（环形）的大小
func SetWindowSize(size int) SliderWindowOption {
	return func(w *SliderWindow) {
		w.size = size
	}
}

// 设置单个窗口代表多长的duration
func SetInterval(interval time.Duration) SliderWindowOption {
	return func(w *SliderWindow) {
		w.interval = interval
	}
}

func SetLogger(logger *zap.Logger) SliderWindowOption {
	return func(w *SliderWindow) {
		w.logger = logger
	}
}
