# Go Context 上下文之源码

## `WithCancel` 的 `cancelCtx` 包含哪些东西？取消是怎么实现的？

1. done：原子创建的 chan struct{}，在调用 Done() 时 make
2. children：存储 `cancelCtx` 和 `timerCtx` 的抽象，每创建一个新 context，就把自己加进来
3. cancel()：关闭 c.done 和 c.children

```go
// A cancelCtx can be canceled. When canceled, it also cancels any children
// that implement canceler.
type cancelCtx struct {
	Context

	mu       sync.Mutex            // protects following fields
	done     atomic.Value          // of chan struct{}, created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      error                 // set to non-nil by the first cancel call
	cause    error                 // set to non-nil by the first cancel call
}
```

## channel 不是并发安全的吗？为什么 `cancelCtx.done` 要用原子操作保证并发安全？

channel 的读写操作是并发安全的，但是获取它的懒汉式单例不是。这跟对象的创建有关，跟 channel 创建后的并发安全没有关系。为了避免被创建多次，就需要用原子操作。
