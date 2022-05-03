package main

import (
	"sync"
	"time"

	"github.com/pandaychen/goes-wrapper/pytime"
	"go.uber.org/zap"
)

// 滑动窗口定义
type SliderWindow struct {
	sync.RWMutex
	offset   int           //记录本次的偏移
	size     int           //大小
	window   *Window       //实际的存储地址
	interval time.Duration //一个窗口对应的duration

	ignoreCurrent bool
	lastTime      time.Duration // start time of the last bucket	（9504h0m0.000065241s）

	logger *zap.Logger
}

// 创建滑动窗口
func NewSliderWindow(opts ...SliderWindowOption) *SliderWindow {
	var (
		w *SliderWindow = new(SliderWindow)
	)

	for _, opt := range opts {
		opt(w)
	}

	if w.size < 1 {
		w.size = 1
	}

	if w.interval == time.Duration(0) {
		w.interval = time.Millisecond * 500
	}

	//set last time
	w.lastTime = pytime.Duration2Now()

	//初始化
	w.window = newWindow(w.size)

	return w
}

// 获取当前时间距离上一次时间之间的窗口数目（计算lastTime 即上一个更新时间到现在跨越了几个 bucket）
func (w *SliderWindow) getSpan() int {
	var (
		// 计算当前时间---lastTime之间的duration，除以w.interval就是滑过了几个窗口
		deltaDuration time.Duration = pytime.Duration2Fixed(w.lastTime)
		offset        int
	)

	// 滑过了多少窗口
	offset = int(deltaDuration / w.interval)

	if 0 <= offset && offset < w.size {
		return offset
	}

	//可能很久都没有滑动了
	return w.size
}

// 执行窗口的滑动（与时间相关）
func (w *SliderWindow) updateOffset() {
	var (
		oldOffset int
		span      int
		nowTime   time.Duration
	)
	span = w.getSpan()
	if span <= 0 {
		//说明当前时间还落在当前的窗口，无须滑动
		return
	}

	oldOffset = w.offset //上一次记录的offset开始
	// 重置过期的窗口
	for i := 0; i < span; i++ {
		w.window.resetBucket((oldOffset + i + 1) % w.size)
	}

	//更新offset
	w.offset = (oldOffset + span) % w.size
	nowTime = pytime.Duration2Now()
	// align to interval time boundary

	//更新lastTime
	w.lastTime = nowTime
}

// 向当前时间对应的bucket中追加数值
// 通过 lastTime 和 nowTime 的不断更新，通过不断重置来实现窗口滑动，新的数据不断补上，从而实现滑动窗口的累加计算
func (w *SliderWindow) Add(v float64) {
	w.Lock()
	defer w.Unlock()
	//先滑动
	w.updateOffset()
	//再加入指标（w.offset是当前时间对应的offset偏移）
	w.window.add(w.offset, v)
}

func main() {
	NewSliderWindow(SetWindowSize(10), SetInterval(500*time.Millisecond))
}
