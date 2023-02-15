# 如何确定 Channel 的缓冲区大小？

缓冲大小为 0 是同步通信，1 则是异步通信。

Uber 的 Go Style Guide 中建议，channel 的缓冲大小要么是 1，要么是无缓冲的，不应该用大于 1 的缓冲。（原因可以参考我的[这篇](https://kanlac.in/channel-buffer-size)博客）

但也有一些特殊场景适合使用大于 1 的缓冲，比如通过 backpressure 机制控制执行某一任务的 goroutine 的数量上限。