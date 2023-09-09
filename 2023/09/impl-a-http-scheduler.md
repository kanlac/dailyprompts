# 实现一个 HTTP 请求调度器

## 题目

假设你正在编写一个支持高并发的 HTTP 请求代理服务。该服务会接收来自多个客户端的 HTTP 请求，并将它们转发到一个内部服务集群。内部服务集群的并发能力有限，所以需要进行请求调度。

**要求：**

1. 实现一个 HTTP 请求调度器，它接收来自客户端的 HTTP 请求。
2. 调度器有一个内部队列，存储待转发的请求。
3. 使用一个池子（数量为 N）的工作 goroutine 来从队列中取出请求并转发到内部服务。
4. 实现请求优先级：有些请求比其他请求更重要，应优先处理。
5. 实现一个超时机制，放弃在队列中等待时间过长的请求。

**考察点：**

1. 使用 channel 和 goroutine 的高级特性。
2. 正确处理并发和同步。
3. 考虑性能和资源利用。
4. 代码组织和模块化。
5. 错误处理和资源清理。

```go
type Request struct {
    ID        int
    Priority  int
    Timestamp time.Time // 请求接收时间
    // 其他字段
}

type Scheduler struct {
    // 你的代码
}

func NewScheduler(workerNum int) *Scheduler {
    // 初始化并返回一个 Scheduler
}

func (s *Scheduler) Enqueue(req Request) {
    // 将请求加入队列
}

func (s *Scheduler) Start() {
    // 启动调度器
}

func (s *Scheduler) Stop() {
    // 停止调度器并清理资源
}

```

**陷阱和挑战：**

1. 如何确保高优先级的请求先被处理？
2. 如何有效地管理工作 goroutine 池？
3. 如何避免队列和工作 goroutine 的死锁？
4. 如何处理超时请求？
5. 如何确保在调度器停止时，所有资源都被正确回收？

这个题目不仅考察基础的并发编程能力，还要求面试者能够进行复杂的系统设计和代码组织。同时，它也有多个维度和方向可以进行扩展和深入，如添加更多的调度算法、错误重试机制等。

## 参考答案

- worker 从队列里拿请求还是直接从 channel 拿请求？——前者的好处是可以做优先级处理，后者的好处是不需要锁

```go
package main

import (
	"fmt"
	"time"
	"container/heap"
	"sync"
)

type Request struct {
	ID        int
	Priority  int
	Timestamp time.Time
}

type PriorityQueue []*Request

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority || (pq[i].Priority == pq[j].Priority && pq[i].Timestamp.Before(pq[j].Timestamp))
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Request)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

type Scheduler struct {
	incomingRequests chan Request
	priorityQueue    PriorityQueue
	lock             sync.Mutex
	workers          int
	timeout          time.Duration
	stopCh           chan bool
}

func NewScheduler(workers int, timeout time.Duration) *Scheduler {
	return &Scheduler{
		incomingRequests: make(chan Request, 100),
		priorityQueue:    make(PriorityQueue, 0),
		workers:          workers,
		timeout:          timeout,
		stopCh:           make(chan bool),
	}
}

func (s *Scheduler) Enqueue(req Request) {
	s.incomingRequests <- req
}

func (s *Scheduler) manageQueue() {
	for {
		select {
		case <-s.stopCh:
			return
		case req := <-s.incomingRequests:
			s.lock.Lock()
			heap.Push(&s.priorityQueue, &req)
			s.lock.Unlock()
		}
	}
}

func (s *Scheduler) worker() {
	for {
		select {
		case <-s.stopCh:
			return
		default:
			s.lock.Lock()
			if s.priorityQueue.Len() == 0 {
				s.lock.Unlock()
				time.Sleep(100 * time.Millisecond)
				continue
			}

			now := time.Now()
			req := heap.Pop(&s.priorityQueue).(*Request)
			s.lock.Unlock()

			if now.Sub(req.Timestamp) > s.timeout {
				fmt.Printf("Request %d timed out\n", req.ID)
				continue
			}
			
			// Process the request here
			fmt.Printf("Processing request %d with priority %d\n", req.ID, req.Priority)
		}
	}
}

func (s *Scheduler) Start() {
	go s.manageQueue()
	for i := 0; i < s.workers; i++ {
		go s.worker()
	}
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}

func main() {
	scheduler := NewScheduler(5, 2*time.Second)
	scheduler.Start()

	// For demo purposes, enqueue some requests
	for i := 0; i < 10; i++ {
		scheduler.Enqueue(Request{ID: i, Priority: i % 3, Timestamp: time.Now()})
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(5 * time.Second)
	scheduler.Stop()
}
```
