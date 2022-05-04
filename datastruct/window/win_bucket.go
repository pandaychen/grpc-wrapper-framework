package window

type WinBucket struct {
	Sum   float64 //成功统计的总数（当前窗口）
	Count int64   //全部（总数，当前窗口）
}

// 累加数值
func (b *WinBucket) add(v float64) {
	b.Sum += v
	b.Count++
}

func (b *WinBucket) reset() {
	b.Sum = 0
	b.Count = 0
}
