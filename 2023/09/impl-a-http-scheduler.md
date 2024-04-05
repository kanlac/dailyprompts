# 实现一个 HTTP 请求调度器

## 正确地提问

1. 希望解决什么问题？
2. 要做的东西是什么？用在什么地方？是服务端的请求处理，还是作为一个反向代理？
3. 是处理同步还是异步请求？
4. 超时如何处理？是否需要重试？重试失败又如何处理？
5. 性能瓶颈在哪里？是否需要做出相应优化？
6. 超过负载上限如何处理？

## 设计与实现

确定要做的东西是什么之后，继续提问，可以选择一个维度深入（比如高并发或者重试机制），得出开发方案，并可以在之后的过程中反复确认、完善和修改方案。

## 作业一：高并发的反向代理

需求如下：编写一个支持高并发的 HTTP 请求代理服务，该服务会接收来自多个客户端的 HTTP 请求，并将它们转发到一个内部服务集群。内部服务集群的并发能力有限，所以需要进行请求调度。功能性要求：

1. 实现 goroutine 的复用
2. 丢弃超时 1 分钟的请求，并返回错误

### 答题框架——使用中间件

main.go

```go
package main

import (
    "log"
    "net/http"
    "time"
)

func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 在这里添加中间件的具体实现

        next.ServeHTTP(w, r)
    })
}

var mockHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
    time.Sleep(200 * time.Millisecond)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("This is a mocked response."))
}

func main() {
    http.Handle("/with-middleware", middleware(mockHandler))
    http.Handle("/without-middleware", mockHandler)

    log.Println("Starting mock proxy server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

```

main_test.go

```go
package main

import (
    "net/http"
    "net/http/httptest"
    "sync"
    "testing"
    "time"
)

func createTestServer() *http.ServeMux {
    mux := http.NewServeMux()
    mux.Handle("/without-middleware", mockHandler)
    mux.Handle("/with-middleware", middleware(mockHandler))
    return mux
}

func TestPerformance(t *testing.T) {
    server := httptest.NewServer(createTestServer())
    defer server.Close()

    var (
        requests   = 1000
        concurrent = 50
    )

    tests := []struct {
        name  string
        route string
    }{
        {"With Middleware", "/with-middleware"},
        {"Without Middleware", "/without-middleware"},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            var wg sync.WaitGroup
            wg.Add(requests)
            start := time.Now()

            for i := 0; i < concurrent; i++ {
                go func() {
                    for j := 0; j < requests/concurrent; j++ {
                        resp, err := http.Get(server.URL + tc.route)
                        if err != nil {
                            t.Errorf("Failed to send request: %v", err)
                        }
                        resp.Body.Close()
                        wg.Done()
                    }
                }()
            }

            wg.Wait()
            elapsed := time.Since(start)
            t.Logf("%s: %d requests took %s", tc.name, requests, elapsed)
        })
    }
}

```

### 验证标准

执行测试用例，处理同样数量的请求，使用中间件消耗时间更长（以此说明转发阻塞生效了）；继续加大请求量，会看到有请求被丢弃

### 参考答案

main.go

```go
package main

import (
"log"
"net/http"
"time"
)

func createLimiterMiddleware(maxConcurrentRequests int, requestTimeout time.Duration) func(http.Handler) http.Handler {
// 使用带缓冲的通道来控制并发量。
semaphore := make(chan struct{}, maxConcurrentRequests)

return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        select {
        case semaphore <- struct{}{}:
            // 请求开始处理，确保释放信号量。
            defer func() { <-semaphore }()
            next.ServeHTTP(w, r)
        case <-time.After(requestTimeout):
            // 如果在requestTimeout时间内没有开始处理请求，则返回超时响应。
            http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
        }
    })
}
}

var mockHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
time.Sleep(200 * time.Millisecond)
w.WriteHeader(http.StatusOK)
w.Write([]byte("This is a mocked response."))
}

func main() {
middleware := createLimiterMiddleware(10, 1*time.Minute)

http.Handle("/with-middleware", middleware(mockHandler))
http.Handle("/without-middleware", mockHandler)

log.Println("Starting mock proxy server on :8080")
if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal("ListenAndServe: ", err)
}
}

```