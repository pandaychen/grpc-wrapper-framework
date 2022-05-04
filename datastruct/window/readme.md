##  window（滑动窗口）

常用于最近一段时间内的行为数据（状态）统计，如自适应熔断器，限流器等；类似于一个简单版本的**时间序列内存数据库**

-   `winBucket`：单个 <窗口>
-   `Window`：环形窗口队列
-   `SliderWindow`：基于 `Window` 基础实现了滑动机制

##  实现思路
-   通过预先创建 `winBucket` 的机制，即当需要记录某个数据（指标）时再去检查并创建 `winBucket`
-   当需要记录指标时，获取当前时间并通过一个算法确定到 `winBucket` ，并记录请求状态


##  WinBucket
```golang
type WinBucket struct {
	Sum   float64 // 成功统计的总数（当前窗口）
	Count int64   // 全部（总数，当前窗口）
}
```

每一个桶 `WinBucket` 存储了：
-   `Sum`：指标总数（通常为 succ）
-   `Count`：请求总数

在 `SliderWindow` 最后聚合计算的时候，通常会将 `Sum` 累计加作为分子，`Count` 累计加作为分母，从而可以统计出当前的成功 / 失败率。

##  窗口滑动的过程

####    滑动的本质
