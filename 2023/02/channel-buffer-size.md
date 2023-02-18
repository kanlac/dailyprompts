# 如何确定 Channel 的缓冲区大小？

缓冲大小为 0 是同步通信，1 则是异步通信。

Uber 的 Go Style Guide 中建议，channel 的缓冲大小要么是 1，要么是无缓冲的，不应该用大于 1 的缓冲。（原因可以参考我的[这篇](https://kanlac.in/channel-buffer-size)博客）

在其它一些特殊场景下，可以使用大于 1 的缓冲：
1. 要发送的数据量是有限的切可以计算出来的。创建 channel 时，使缓冲区大小与输入的数据量相同，这样一来，就不需要开新协程了
2. 通过 Backpressure 机制控制执行某一任务的 goroutine 的数量上限