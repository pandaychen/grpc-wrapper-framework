package window

type WinBucket struct {
	Sum   float64 // 成功统计的总数（当前窗口）
	Count int64   // 全部（总数，当前窗口）
}

// 向桶添加数据，累加数值
func (b *WinBucket) add(v float64) {
	b.Sum += v // 累加指标
	b.Count++  // 累加总次数
}

// 桶数据清零
func (b *WinBucket) reset() {
	b.Sum = 0
	b.Count = 0
}
