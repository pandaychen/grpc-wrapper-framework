package window

import (
	//"math/rand"
	"fmt"
	"sync"
	"time"

	"github.com/pandaychen/goes-wrapper/pytime"
	"go.uber.org/zap"
)

// 滑动窗口定义
type SliderWindow struct {
	sync.RWMutex
	offset   int           //记录本次的偏移
	size     int           //大小（可能需要累加到size？）
	window   *Window       //实际的存储地址
	interval time.Duration //一个窗口对应的duration

	ignoreCurrentBucket bool          //统计的时候是否忽略当前bucket（数据不全）
	lastTime            time.Duration // start time of the last bucket      （9504h0m0.000065241s）

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

// 聚合滑动窗口中的所有有效bucket的数据
func (w *SliderWindow) Reduce(fn func(b *WinBucket)) {
	var (
		span  int
		start int
		count int
	)
	w.RLock()
	defer w.RUnlock()

	// 计算当前时间截止前，已过期桶的数量
	span = w.getSpan()
	// ignore current bucket, because of partial data
	//这里span非0的话，说明当前时间离最近一次时间有跨度，需要排除掉（这些跨度中无新数据）
	if span == 0 && w.ignoreCurrentBucket {
		count = w.size - 1
	} else {
		count = w.size - span
	}
	if count > 0 {
		// w.offset 与 rw.offset+span之间的桶数据是过期的，不应该计入统计
		start = (w.offset + span + 1) % w.size //这里暴力的轮询完整个size的窗口数目

		//聚合数据
		w.window.reduce(start, count, fn)
	}
}

func main() {
	var (
		winSize  = 4
		duration = time.Millisecond * 500
	)
	window := NewSliderWindow(SetWindowSize(winSize), SetInterval(duration))
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				window.Add(float64(2))
				time.Sleep(duration / time.Duration(winSize))
			}
		}
	}()
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				var (
					result float64
					count  int64
				)
				time.Sleep(time.Duration(winSize) * duration)
				window.Reduce(func(b *WinBucket) {
					result += b.Sum
					count += b.Count

				})
				fmt.Println("result=", result, "count=", count)
			}
		}
	}()
	time.Sleep(duration * 5000)
	close(stop)
}
