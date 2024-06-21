# Go Context 上下文之源码

## 从源码层面解读上下文的实现原理

1. WithCancel(), WithDeadline(), WithTimeout(), WithValue()
2. `Context` 接口：`Done()` 获取终止信号 channel，`Err()` 获取错误，比如取消或超时
3. `cancelCtx` 包含 `Context` 接口
    - `children map[canceler]struct{}` - 存储 cancelCtx
    - `done atomic.Value` - **懒汉**式创建的 `chan struct{}`，调用 Done() 时返回该值
4. `timerCtx`，包含 `cancelCtx`，`deadline time.Time`
5. `valueCtx`, 包含 `Context` 接口，`key, val any`

## channel 不是并发安全的吗？为什么 cancelCtx.done 要用原子操作保证并发安全？

channel 的读写操作是并发安全的，但是获取它的懒汉式单例不是。这跟对象的创建有关，跟 channel 创建后的并发安全没有关系。为了避免被创建多次，就需要用原子操作。
    