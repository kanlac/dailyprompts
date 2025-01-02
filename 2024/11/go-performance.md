# Go 性能分析

## 性能分析有哪几种数据（Types of profiles）

1. CPU profile：帮助定位使用较多 CPU 时间的函数
2. Memory profile：帮助定位应用中哪部分使用了较多内存
3. Block profile：帮助定位程序在哪个地方阻塞时间较长，优化并发
4. Goroutine profile：通过运行、阻塞、等待等状态，帮助定位哪里存在较多并行，优化并发
5. Trace profile：提供事件日志，比如 goroutine 创建和销毁、调度、网络活动和阻塞操作等，用于更细致的分析

## 如何使用 pprof

Go 提供了内建的分析工具，pprof，它可以用来分析程序的内存占用情况。pprof 可以生成堆内存的剖析文件，然后你可以使用 `go tool pprof` 命令来分析这个文件。

使用 pprof 的步骤：

- 在代码中导入 `net/http/pprof`。
- 启动一个 HTTP 服务器来提供剖析接口。
- 使用 `go tool pprof` 工具来获取和分析数据。

示例代码：

```go
import (
    _ "net/http/pprof"
    "net/http"
)

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    // 其他代码
}
```

然后，在终端中运行 `go tool pprof <http://localhost:6060/debug/pprof/heap`。

## pprof 原理是什么

在指定周期内，按照频率（默认是 10ms）采样调用栈信息。这样，某个调用栈的汇总 CPU 时间就等于该调用栈的采样次数 * 采样间隔。采样频率太低了不精确，太高了可能会影响性能。

如何指定间隔周期？——以 http hander 这种 profile 方式为例，就是每次访问特定 route 的时候采样 30s，源码这里会 sleep，然后再显示结果。

pprof 对性能影响很小，是可以在生产环境下使用的。

再往下看有点细节了……参见文章吧。

## CPU profile 数据怎么读

pprof 采样数据显示某个函数节点 0.01s (0.0086%) of 17.56s (15.03%)，这是什么意思？——前者是该函数本身的时间，很少；后者是该函数及其它所调用的所有函数的总时间，占比有点高。意思就是要关注这个函数或某个自定义的子函数。

## 性能分析优化实践

[万字长文讲透Go程序性能优化](https://mp.weixin.qq.com/s/wLPfiJ0wKH3DrBJS4yxeHw)

1. debug 日志打印，本来不应该打印但产生了耗时，优化方式是把参数中的 String() 方法改为传递指针
2. 不合理的库函数调用，比如 runtime.growslice 反应切片的自动扩展，可以通过初始化时指定长度降低
3. gc 相关，runtime.gcBgMarkWorker 占比高反应 GC 频率较高，可以通过调**高** GOGC 优化