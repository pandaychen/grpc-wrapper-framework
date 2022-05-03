package main

// window is  like a ring queue
type Window struct {
	buckets []*WinBucket
	size    int //ring
}

//
func newWindow(size int) *Window {
	buckets := make([]*WinBucket, size)

	//init
	for i := 0; i < size; i++ {
		buckets[i] = new(WinBucket)
	}
	return &Window{
		buckets: buckets,
		size:    size,
	}
}

func (w *Window) add(offset int, v float64) {
	// 向offset指向的 bucket 加入指定的指标
	if bucket := w.buckets[offset%w.size]; bucket != nil {
		bucket.add(v)
	}
}

// reset window bucket
func (w *Window) resetBucket(offset int) {
	if bucket := w.buckets[offset%w.size]; bucket != nil {
		bucket.reset()
	}
}
