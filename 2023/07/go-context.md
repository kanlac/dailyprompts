# Go Context 上下文

## 如果没有上下文，同类问题如何解决？有什么问题？说出两个具体场景
    
`context` 这个包到底解决了什么问题？很简单，就是如何取消（cancel）协程。在 `context` 包出来之前，有过一个讲座专门谈论协程的取消问题，解决方案用原生 channel 实现，但伸缩性较差。主要问题是：

1. 别的库不会接收 cancel channel，所以只能在不同的操作之间做取消
2. 设想一个 goroutine 树，要直接取消整个树是很简单的，但如果要取消一个子树，你还需要定义一个新的 cancel channel

随后 `context` 包的出现，解决了这两个问题，尽管它存在问题，但目前仍是最好的方案。

## 从源码层面解读上下文的实现原理

1. WithCancel(), WithDeadline(), WithTimeout(), WithValue()
2. `Context` 接口：`Done()` 获取终止信号 channel，`Err()` 获取错误，比如取消或超时
3. `cancelCtx` 包含 `Context` 接口
    - `children map[canceler]struct{}` - 存储 cancelCtx
    - `done atomic.Value` - 懒汉式创建的 `chan struct{}`，调用 Done() 时返回该值
4. `timerCtx`，包含 `cancelCtx`，`deadline time.Time`
5. `valueCtx`, 包含 `Context` 接口，`key, val any`
    
## 为什么说给 `io.Reader` 加上上下文不是一个好提议？

1. Go 本身定位应该是通用型语言（general purpose language），而不是服务端专用语言，而上下文基本都是写服务的时候使用
2. 上下文是病毒
3. 代码不优雅

## HTTP 请求关闭时 Go 服务端的上下文是否会取消？有什么特例？

一般来说，关闭页面，刷新页面，或者 Ctrl+C 中止 curl 命令，都会关闭 TCP 连接，这样的话就 http 包的 context 就能正常 cancel。

但是有例外情况。

一）有请求转发，并且转发后的连接不会随着前面的连接关闭而关闭，这样就造成服务端感知不到连接关闭了。

二）没有读完 Body，或者没有 close body 的情况下，即使连接关闭也不会 cancel，这个比较坑，是 Go 实现的一个问题，17 年的 issue 一直悬而未决 https://github.com/golang/go/issues/23262
