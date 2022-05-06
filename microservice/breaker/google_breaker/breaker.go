package google_breaker

//google 熔断算法实现

import (
	"grpc-wrapper-framework/config"
	wd "grpc-wrapper-framework/datastruct/window"
	"grpc-wrapper-framework/errcode"
	bc "grpc-wrapper-framework/microservice/breaker/common"
	"math"
	"sync/atomic"

	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	//
	PROTECTION_LIMIT int64 = 5
)

// googleSREBreaker is a sre CircuitBreaker pattern.
type googleSREBreaker struct {
	name     string
	rnd      *rand.Rand //rand.New() 非线程安全
	randLock sync.Mutex
	kval     float64
	request  int64 // 达到请求上限才触发检测
	state    int32
	// 滑动窗口累计
	statistics *wd.SliderWindow
	loger      *zap.Logger
}

func NewBreaker(logger *zap.Logger, c *config.BreakerConfig) bc.Breaker { //interface{}
	var (
		window *wd.SliderWindow
	)
	//init slider window
	window = wd.NewSliderWindow(
		wd.SetWindowSize(c.GoogleBreakerConf.BucketSize),
		wd.SetInterval(c.GoogleBreakerConf.WindowDuration),
	)
	return &googleSREBreaker{
		rnd:        rand.New(rand.NewSource(time.Now().UnixNano())),
		request:    c.GoogleBreakerConf.Request,
		kval:       c.GoogleBreakerConf.Kval,
		state:      bc.StateClosed,
		statistics: window,
		loger:      logger,
		name:       c.GoogleBreakerConf.Name,
	}
}

func (b *googleSREBreaker) Name() string {
	return b.name
}

// 统计滑动窗口的指标
func (b *googleSREBreaker) summary() (int64, int64) {
	var (
		accepts, total int64
	)

	// 设置统计方法
	b.statistics.Reduce(func(b *wd.WinBucket) {
		accepts += int64(b.Sum)
		total += b.Count
	})

	return accepts, total
}

func (b *googleSREBreaker) MarkSuccess() {
	// 成功累加
	b.statistics.Add(1)
}

func (b *googleSREBreaker) MarkFailed() {
	// 失败累加，强制窗口滑动以更新计算错误率
	b.statistics.Add(0)
}

// 返回熔断器当前状态
func (b *googleSREBreaker) State() int32 {
	return atomic.LoadInt32(&b.state)
}

func (b *googleSREBreaker) trueOnProbability(proba float64) bool {
	var (
		truth bool
	)
	b.randLock.Lock()
	truth = b.rnd.Float64() < proba
	b.randLock.Unlock()
	return truth
}

// 实时：当前熔断器状态是否允许请求放过
func (b *googleSREBreaker) Allow() error {
	var (
		succ, total int64
		tolerance   float64
		// 客户端请求拒绝的概率
		dropRatio float64
	)
	// 获取滑动窗口的瞬时统计值
	succ, total = b.summary()
	tolerance = b.kval * float64(succ) //K*success
	if total < b.request || float64(total) < tolerance {
		// 放行
		return nil
	}

	if dropRatio = math.Max(0, (float64(total-PROTECTION_LIMIT)-tolerance)/float64(total+1)); dropRatio <= 0 {
		// 放行
		return nil
	}

	//drop with fixed probability
	if b.trueOnProbability(dropRatio) {
		// 返回全局错误码，丢弃请求
		return errcode.ServiceUnavailable
	}

	// 放行
	return nil
}

// 用户方法调用
func (b *googleSREBreaker) DoRequest(usercall func() error, fallback bc.FallBackCaller, isErrAccept bc.CheckIfAcceptableCaller) error {
	var (
		err error
	)
	if err = b.Allow(); err != nil {
		if fallback != nil {
			// 熔断状态开启，执行用户传入的 fallback 方法
			return fallback(err)
		}
	}

	//Do real usercall
	if usercall != nil {
		err = usercall()
		if isErrAccept != nil {
			if isErrAccept(err) {
				// 非熔断所触发的错误类型，如逻辑错误等
				b.MarkSuccess()
			} else {
				b.MarkFailed()
			}
		} else {
			if err != nil {
				b.MarkFailed()
			} else {
				b.MarkSuccess()
			}
		}
	}

	return err
}
