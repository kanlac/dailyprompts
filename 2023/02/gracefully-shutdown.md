# 如何实现 Go 程序的优雅退出

### 基本概念
优雅退出对于长期运行的有状态程序很有意义，它是 SRE 站点可靠性工程中的一个重要方面。

优雅退出就是要有一个机制，使得上层协程能够在监听到中止信号时通知下层协程，这样在机器关机，或者 Kubernetes pod 结束时，程序能完成一些有必要的操作。必要操作包括：
- 在关闭服务器前，完成接收到退出信号之前发起的请求，并不再接受新的请求
- 在关闭数据库前，保存当前状态
- 关闭所有外部服务和数据库的连接

### 反设计模式/不要怎么做
一）不要持续阻塞 main 进程不释放：
```go
func KeepProcessAlive() {
	var ch chan int
	<-ch
}

func main() {

	//...
	KeepProcessAlive()
}
```

二）不要使用 `os.Exit(1)`，这样就没法提前关闭开启的连接。

几种常用 UNIX shutdown 信号：
- SIGTERM｜大部分 shutdown 事件都发送这种终止信号
- SIGKILL｜「立即退出」，最好不要用（避免使用 `os.Exit(1)`）
- SIGINT｜和 SIGTERM 类似，但代表的是用户输入的中断信号（如 ctrl + C）
- SIGQUIT｜和 SIGKILL 类似，但代表的是用户输入的退出信号（如 ctrl + D）

生产环境下，监听 SIGTERM 足够了（Kubernetes pods 结束时发送的就是这个信号）；SIGINT 则可以用于本地调试。处理好这两个信号可以覆盖绝大多数场景。

### 如何做
两个目标/两个问题：

1. 如何等待所有进行中的 goroutines 退出——使用 channel 或者 waitgroup 都行，后者更简洁

2. 如何传递终止信号到各个 goroutines

创建一个 channel，监听退出信号：

```go
func main() {
	gracefulShutdown := make(chan os.Signal, 1)  // 为了频道的发送方不阻塞，设置缓冲区大小为 1
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)
	// do stuff
	// ...
	<-gracefulShutdown
	// do stuff
	fmt.Println("Done!")
}
```

当 gracefulShutdown 监听到了 SIGINT 或者 SIGTERM 信号时，就会打印上面的 “Done!”。

如果只有一个 goroutine 需要优雅退出，那么使它监听这个 `gracefulShutdown` channel 就足够了。但如果有多个呢？毕竟 channel 只能读一次。这时候就要用到 Context 上下文：

```go
func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	// watch os signal
	go func() {
		<-exit
		fmt.Println("\nreceive terminate signal, start graceful stop process")
		cancel()
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("gracefully stop process one")
				return
			default:
				fmt.Println("running process one")
				time.Sleep(time.Second * 3)
			}
		}
	}(ctx)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("gracefully stop process two")
				return
			default:
				fmt.Println("running process two")
				time.Sleep(time.Second * 3)
			}
		}
	}(ctx)

	wg.Wait()
	fmt.Println("done!")
}

/* output:
running process two
running process one
running process one
running process two
running process two
running process one
^C
receive terminate signal, start graceful stop process
gracefully stop process two
gracefully stop process one
done!
*/
```