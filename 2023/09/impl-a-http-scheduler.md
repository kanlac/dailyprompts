# 实现一个 HTTP 请求调度器

## 有哪些问题需要注意和明确

问出必要的、好的问题，是展现工程师素养的第一步。

1. 讨论清楚要做的东西是什么？是作为服务端并发处理请求，还是作为一个反向代理？
2. 基本功能方面，是否需要做超时？超时如何处理？请求是同步处理还是异步处理？
3. 如果要求在某个维度上深入（比如实现重试机制），再进一步询问……

## 编写调用端，进一步明确需求

在需求初步明确之后，编写调用端代码，询问是否符合预期，并作出调整和修改。

现假设需求如下：编写一个支持高并发的 HTTP 请求代理服务，该服务会接收来自多个客户端的 HTTP 请求，并将它们转发到一个内部服务集群。内部服务集群的并发能力有限，所以需要进行请求调度。功能性要求：

1. 实现 goroutine 的复用
2. 丢弃超时 1 分钟的请求，并返回错误

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// server
	scheduler := NewScheduler(5, 1*time.Second)
	scheduler.Start() // listen on 80
	defer scheduler.Stop()

	// client
	errCh := make(chan error)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go makeHTTPRequest(errCh, &wg)
	}

	go func() {
		wg.Wait()
		fmt.Println("all request finished")
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			fmt.Printf("err: %+v\n", err)
		}
	}
}
```

## 编写实现

……
