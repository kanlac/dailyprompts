# Golang 协程池技术调研

## What is the purpose of a Goroutine Pool in Go, and how can it be used to manage concurrent tasks efficiently?

A Goroutine Pool, also known as a worker pool or task pool, is a design pattern that involves pre-creating a fixed number of Goroutines to handle tasks concurrently. The purpose of a Goroutine Pool is to limit the number of concurrent Goroutines and efficiently manage task execution. By reusing Goroutines, you can avoid the overhead of creating and destroying Goroutines for each task, resulting in better performance and resource management.

In Go, a Goroutine Pool can be implemented using a combination of Goroutines, channels, and synchronization mechanisms like WaitGroups or semaphores. Tasks are typically submitted to the pool through a channel, and the Goroutines in the pool continuously read from the channel, execute the tasks, and then loop back to read from the channel again. By controlling the number of Goroutines in the pool and the size of the task queue, you can fine-tune the level of concurrency and prevent resource exhaustion.

[source](https://coderpad.io/interview-questions/go-interview-questions/)

## 实现一个协程池

客户端：

```go
func task(ctx context.Context) error {
		fmt.Println("处理任务", count)
		time.Sleep(time.Second)
}

func main() {
    pool := NewGoroutinePool(4, 10) // 创建一个大小为 4 的协程池

    for i := 0; i < 10; i++ {
        count := i
        pool.Go(func(ctx context.Context) error {
						return task(ctx)
        })
    }

    err := pool.Wait() // 阻塞，等待所有任务完成，并收集错误
}
```

## 协程池技术调研

### 仓库地址与社区

- `errgroup` 官方库 [golang.org/x/sync/errgroup](http://golang.org/x/sync/errgroup)
- `ants`  https://github.com/panjf2000/ants, 9.9k star
- `conc` https://github.com/sourcegraph/conc, 6k star

### 太长不看

- 使用协程池的目的：1）使并发程序代码更简洁；2）高并发情况下，节省内存占用
- 特性支持情况：

|  | 池中 worker 可复用 | 上下文 | WaitGroup 封装 | 错误处理 | 动态调整池容量 |
| --- | --- | --- | --- | --- | --- |
| errgroup | ❌ | ✅ | ✅ | ✅ | ❌ |
| conc | ✅ | ✅ | ✅ | ✅ | ❌ |
| ants | ✅ | ❌ | ❌ | ❌ | ✅ |
- `conc` 的功能最全面，基本可以覆盖各种场景，但因为使用泛型特性，需要 Go 1.19+
- `conc` 对 panic 的处理：`conc` 有一个 re-panic 的机制，它会先捕获 task 中产生的 panic，等待资源回收完毕，最后 `Wait()` 再重新产生 panic。因此，可以按需在 `Wait()` 的调用处自行 recover
- 不考虑风格偏好，从使用简便性来看，`errgroup` 和 `conc` 都相当出色，`ants` 则比较逊色
- `conc` 的内存优化最佳，`ants` 其次（见后文性能测试）

### 源码分析

- `ants` 维护一个 worker 队列，每一个 worker 有一个任务 channel
- `conc` 维护一个自己的 worker 队列，整个协程池只有一个任务 channel
- 相对的，用户每次传入一个任务，`errgroup` 都会创建一个新的协程

关于 channel 的缓冲区大小

- `ants` 在创建 worker 的同时创建新 channel，若 GOMAXPROCS=1，则采用同步通信（buffer=0），若 GOMAXPROCS>1，则采用异步通信（buffer=1），据说性能更优（详情参见注释）
- `conc` 采用同步通信（buffer=0）

### 使用演示代码

场景设例：计算目录下文件的 MD5 值。

- demo
    
    ```go
    package main
    
    import (
    	"context"
    	"crypto/md5"
    	"fmt"
    	"io/ioutil"
    	"log"
    	"os"
    	"path/filepath"
    	"sync"
    
    	"github.com/panjf2000/ants"
    	"github.com/sourcegraph/conc/pool"
    	"golang.org/x/sync/errgroup"
    )
    
    // Pipeline demonstrates the use of a Group to implement a multi-stage
    // pipeline: a version of the MD5All function with bounded parallelism from
    // https://blog.golang.org/pipelines.
    func main() {
    	m, err := MD5All(context.Background(), ".")
    	if err != nil {
    		log.Fatal(err)
    	}
    
    	for k, sum := range m {
    		fmt.Printf("%s:\t%x\n", k, sum)
    	}
    }
    
    type result struct {
    	path string
    	sum  [md5.Size]byte
    }
    
    // MD5All reads all the files in the file tree rooted at root and returns a map
    // from file path to the MD5 sum of the file's contents. If the directory walk
    // fails or any read operation fails, MD5All returns an error.
    func MD5All(ctx context.Context, root string) (map[string][md5.Size]byte, error) {
    	paths := make(chan string)
    
    	var walkFunc filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
    		if err != nil {
    			return err
    		}
    		if !info.Mode().IsRegular() {
    			return nil
    		}
    		select {
    		case paths <- path:
    		case <-ctx.Done():
    			return ctx.Err()
    		}
    		return nil
    	}
    
    	const poolSize = 1000
    
    	// ======================== 协程池的创建 ========================
    
    	// >>>>>> errgroup
    	// 支持上下文，当 task 发生错误时，errgroup 将取消上下文
    	// errgroup 不支持通过 task 传参，因此只能通过 channel传递
    	errgroupPool, ctx := errgroup.WithContext(ctx)
    	errgroupPool.SetLimit(poolSize)
    
    	// >>>>>> conc
    	concPool := pool.New().
    		WithContext(ctx).
    		WithMaxGoroutines(poolSize).
    		WithCancelOnError().
    		WithFirstError()
    
    	// >>>>>> ants
    	// 选择一：通用池（NewPool），task signature: func()
    	// 选择二：专用池（NewPoolWithFunc），task signature: func(interface{})
    	// 专用池创建时定义 task func，可以简化调用代码
    	// 缺点：创建协程池时需要考虑通用池和专用池哪种合适
    	// 使用专用池看似可以简化调用的代码，但实际上池的创建更麻烦（要做 assert）
    	// 缺点：不支持上下文
    	// 优势：支持通过 Tune() 动态、且线程安全地调整协程池大小（errgroup 需要等协程池空了之后再调整，conc 无法调整）
    	// 因为 stage 1 只需要开一个 worker，所以实际上专门创建一个协程池意义不大，这里只出于演示目的创建一个专用池
    	antsPoolS1, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
    		root := i.(string)
    		_ = filepath.Walk(root, walkFunc) // 错误捕获需要手动完成（例如发到 channel）
    		close(paths)
    	})
    	defer antsPoolS1.Release()
    
    	// ======================== Stage 1 ========================
    
    	// >>>>>> errgroup
    	// 支持且仅支持的 task signature: func() error（足够用）
    	errgroupPool.Go(func() error {
    		defer close(paths)
    		return filepath.Walk(root, walkFunc)
    	})
    
    	// >>>>>> conc
    	concPool.Go(func(ctx context.Context) error {
    		// 可以在这里拿到上下文
    		defer close(paths)
    		return filepath.Walk(root, walkFunc)
    	})
    
    	// >>>>>> ants
    	if err := antsPoolS1.Invoke(root); err != nil {
    		return nil, err
    	}
    
    	// ======================== Stage 2 ========================
    
    	// Start a fixed number of goroutines to read and digest files.
    	c := make(chan result)
    	worker2 := func() error {
    		for path := range paths {
    			data, err := ioutil.ReadFile(path)
    			if err != nil {
    				return err
    			}
    			select {
    			case c <- result{path, md5.Sum(data)}:
    			case <-ctx.Done():
    				return ctx.Err()
    			}
    		}
    		return nil
    	}
    	const stage2WorkerNum = 5
    
    	// >>>>>> errgroup
    	// errgroup 集成了 sync.WaitGroup
    	for i := 0; i < stage2WorkerNum; i++ {
    		errgroupPool.Go(worker2)
    	}
    
    	// >>>>>> conc
    	for i := 0; i < stage2WorkerNum; i++ {
    		concPool.Go(func(ctx context.Context) error {
    			return worker2()
    		})
    	}
    
    	// >>>>>> ants
    	// 因为专用池的创建与 task func 是耦合的，所以用新的 task 要创建新的协程池
    	// 这次演示通用池
    	var wg sync.WaitGroup
    	antsPoolS2, _ := ants.NewPool(poolSize)
    	defer antsPoolS2.Release()
    	for i := 0; i < stage2WorkerNum; i++ {
    		wg.Add(1)
    		err := antsPoolS2.Submit(func() {
    			_ = worker2()
    			wg.Done()
    		})
    		if err != nil {
    			return nil, err
    		}
    	}
    
    	var concPoolErr error
    	go func() {
    		// >>>>>> errgroup
    		errgroupPool.Wait()
    
    		// >>>>>> conc
    		concPoolErr = concPool.Wait()
    
    		// >>>>>> ants
    		wg.Wait()
    		fmt.Printf("running goroutines: %d\n", antsPoolS1.Running())
    
    		// ===
    		close(c)
    	}()
    
    	m := make(map[string][md5.Size]byte)
    	for r := range c {
    		m[r.path] = r.sum
    	}
    
    	// ======================== 错误处理 ========================
    
    	// >>>>>> errgroup
    	// Check whether any of the goroutines failed. Since g is accumulating the
    	// errors, we don't need to send them (or check for them) in the individual
    	// results sent on the channel.
    	if err := errgroupPool.Wait(); err != nil {
    		return nil, err
    	}
    
    	// >>>>>> conc
    	if concPoolErr != nil {
    		return nil, concPoolErr
    	}
    
    	// >>>>>> ants
    	// 不支持
    
    	// ======
    
    	return m, nil
    }
    ```
    

### 性能测试数据

- Apple M1 八核
    
    第一组
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/0f5d1d66-7127-4dc1-8383-066769714e4c/Untitled.png)
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/05e2e160-6375-4e6f-b07f-a312be743de3/Untitled.png)
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/f4aa7c97-8e06-4e87-b87d-fb0d21e8ff00/Untitled.png)
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/d0f2e29a-fbae-4c60-8613-11467fb3f89c/Untitled.png)
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/a3eca69b-2786-4de5-a898-a63cc76baea6/Untitled.png)
    
    第二组 10s
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/8eab2a1c-9aa9-435f-bb69-aecdbcf40267/Untitled.png)
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/e9f1048c-6620-4d79-aa71-4e742762b147/Untitled.png)
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/dd99db0c-50b1-42e1-928d-01aa1fa376dc/Untitled.png)
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/85c36f04-b5ca-4985-8738-13f2be0f58bd/Untitled.png)
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/335cfbb5-5d80-4c8e-86db-ad47ee4bbaae/Untitled.png)
    
    第三组 50s
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/124c33f0-8467-43c5-b33f-26dd855af9d8/Untitled.png)
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/0dacf0d1-bfc0-4267-8c3a-2af836acbefc/Untitled.png)
    
- Apple M1 单核
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/5f13211c-7616-49e2-8777-440600dcbcdc/Untitled.png)
    
- 开发机（双核）
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/eebf4878-0d1d-4935-b057-2ac5cebf6a33/Untitled.png)
    
- 开发机（单核）
    
    ![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/3bad7317-5f71-4557-83c2-eb251defa74c/Untitled.png)
    
- 测试代码
    
    基于 ants_benchmark_test.go 修改
    
    ```go
    // MIT License
    
    // Copyright (c) 2018 Andy Pan
    
    // Permission is hereby granted, free of charge, to any person obtaining a copy
    // of this software and associated documentation files (the "Software"), to deal
    // in the Software without restriction, including without limitation the rights
    // to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    // copies of the Software, and to permit persons to whom the Software is
    // furnished to do so, subject to the following conditions:
    //
    // The above copyright notice and this permission notice shall be included in all
    // copies or substantial portions of the Software.
    //
    // THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    // IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    // FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    // AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    // LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    // OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    // SOFTWARE.
    
    package ants
    
    import (
    	"runtime"
    	"sync"
    	"testing"
    	"time"
    
    	"github.com/sourcegraph/conc/pool"
    	"golang.org/x/sync/errgroup"
    )
    
    const (
    	RunTimes           = 1e6
    	PoolCap            = 5e4
    	BenchParam         = 10
    	DefaultExpiredTime = 10 * time.Second
    )
    
    func demoFunc() {
    	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
    }
    
    func demoPoolFunc(args interface{}) {
    	n := args.(int)
    	time.Sleep(time.Duration(n) * time.Millisecond)
    }
    
    func longRunningFunc() {
    	for {
    		runtime.Gosched()
    	}
    }
    
    func longRunningPoolFunc(arg interface{}) {
    	if ch, ok := arg.(chan struct{}); ok {
    		<-ch
    		return
    	}
    	for {
    		runtime.Gosched()
    	}
    }
    
    func BenchmarkGoroutines(b *testing.B) {
    	var wg sync.WaitGroup
    	for i := 0; i < b.N; i++ {
    		wg.Add(RunTimes)
    		for j := 0; j < RunTimes; j++ {
    			go func() {
    				demoFunc()
    				wg.Done()
    			}()
    		}
    		wg.Wait()
    	}
    }
    
    // func BenchmarkChannel(b *testing.B) {
    // 	var wg sync.WaitGroup
    // 	sema := make(chan struct{}, PoolCap)
    
    // 	b.ResetTimer()
    // 	for i := 0; i < b.N; i++ {
    // 		wg.Add(RunTimes)
    // 		for j := 0; j < RunTimes; j++ {
    // 			sema <- struct{}{}
    // 			go func() {
    // 				demoFunc()
    // 				<-sema
    // 				wg.Done()
    // 			}()
    // 		}
    // 		wg.Wait()
    // 	}
    // }
    
    func BenchmarkErrGroup(b *testing.B) {
    	var wg sync.WaitGroup
    	var pool errgroup.Group
    	pool.SetLimit(PoolCap)
    
    	b.ResetTimer()
    	for i := 0; i < b.N; i++ {
    		wg.Add(RunTimes)
    		for j := 0; j < RunTimes; j++ {
    			pool.Go(func() error {
    				demoFunc()
    				wg.Done()
    				return nil
    			})
    		}
    		wg.Wait()
    	}
    }
    
    func BenchmarkErrGroupCustom(b *testing.B) {
    	var pool errgroup.Group
    	pool.SetLimit(PoolCap)
    
    	b.ResetTimer()
    	for i := 0; i < b.N; i++ {
    		for j := 0; j < RunTimes; j++ {
    			pool.Go(func() error {
    				demoFunc()
    				return nil
    			})
    		}
    		pool.Wait()
    	}
    }
    
    func BenchmarkConcPool(b *testing.B) {
    	pools := make([]*pool.Pool, b.N)
    	for i := 0; i < b.N; i++ {
    		pools[i] = pool.New().WithMaxGoroutines(PoolCap)
    	}
    
    	b.ResetTimer()
    	for i := 0; i < b.N; i++ {
    		for j := 0; j < RunTimes; j++ {
    			pools[i].Go(demoFunc)
    		}
    		pools[i].Wait()
    	}
    }
    
    func BenchmarkAntsPool(b *testing.B) {
    	var wg sync.WaitGroup
    	p, _ := NewPool(PoolCap, WithExpiryDuration(DefaultExpiredTime))
    	defer p.Release()
    
    	b.ResetTimer()
    	for i := 0; i < b.N; i++ {
    		wg.Add(RunTimes)
    		for j := 0; j < RunTimes; j++ {
    			_ = p.Submit(func() {
    				demoFunc()
    				wg.Done()
    			})
    		}
    		wg.Wait()
    	}
    }
    
    // func BenchmarkGoroutinesThroughput(b *testing.B) {
    // 	for i := 0; i < b.N; i++ {
    // 		for j := 0; j < RunTimes; j++ {
    // 			go demoFunc()
    // 		}
    // 	}
    // }
    
    // func BenchmarkSemaphoreThroughput(b *testing.B) {
    // 	sema := make(chan struct{}, PoolCap)
    // 	for i := 0; i < b.N; i++ {
    // 		for j := 0; j < RunTimes; j++ {
    // 			sema <- struct{}{}
    // 			go func() {
    // 				demoFunc()
    // 				<-sema
    // 			}()
    // 		}
    // 	}
    // }
    
    // func BenchmarkAntsPoolThroughput(b *testing.B) {
    // 	p, _ := NewPool(PoolCap, WithExpiryDuration(DefaultExpiredTime))
    // 	defer p.Release()
    
    // 	b.ResetTimer()
    // 	for i := 0; i < b.N; i++ {
    // 		for j := 0; j < RunTimes; j++ {
    // 			_ = p.Submit(demoFunc)
    // 		}
    // 	}
    // }
    ```
    

### 性能数据分析

M1 八核

- 计算性能按排名，纯原生约为 0.5 s/op，errgroup 与 ants 的计算性能均为 0.8-0.9 s/op，conc 约为 1.1 s/op
- 内存占用按排名，conc 约为 6 MB/op，ants 约为 17 MB/op，纯原生约为 96 MB/op，errgroup 约为 104 MB/op

M1 单核

- 计算性能按排名，纯原生和 conc 均为 1.3 s/op，errgroup 1.4 s/op，ants 1.5 s/op，整体差异小
- 内存占用按排名，conc 10 MB/op，ants 32 MB/op，errgroup 104 MB/op，纯原生 96-120 MB/op 波动

虚拟机双核

- 计算性能按排名，ants 约为 2 s/op，纯原生约为 2.1 s/op，errgroup 约为 2.2 s/op，conc 约为 2.5 s/op
- 内存占用按排名，conc 10 MB/op，ants 22 MB/op，errgroup 104 MB/op，纯原生 96-120 MB/op 波动

虚拟机单核

- 计算性能按排名，纯原生 3.4 s/op，errgroup 3.6 s/op，conc 3.9 s/op，ants 4.6 s/op
- 内存占用按排名，conc 10 MB/op，ants 32 MB/op，errgroup 104 MB/op，纯原生 96-120 MB/op 波动

### 性能分析总结

纯原生 goroutine 的计算性能已经无可比拟，协程库的优势主要体现在内存占用上，尤其是 conc 表现最佳，比纯原生实现内存占用少 10 倍。