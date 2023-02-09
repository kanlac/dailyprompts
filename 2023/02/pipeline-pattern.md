### 流水线模式是 Go 中比较常见的一个并发模式，它在什么情况下使用？请给出代码实现

流水线模式模拟的是工厂里的生产流水线，我们可以指定开启一定数量的 goroutine 来处理耗时的计算操作。完成流水线的组装后，它会监听输入源 channel 并自动进行并发处理。

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func GetInputChan() <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func GeneratePipeline(input <-chan int) <-chan int {
	ch := make(chan int)
	go func() {
		for v := range input {
			// do some time consuming work on v
			time.Sleep(3 * time.Second)
			ch <- v * v
		}
		close(ch)
	}()
	return ch
}

func MergePipelines(pipelines ...<-chan int) <-chan int {
	ch := make(chan int)
	var wg sync.WaitGroup
	for _, p := range pipelines {
		wg.Add(1)
		localPipeline := p
		go func() {
			for v := range localPipeline {
				ch <- v
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func PersistResult(ch <-chan int) {
	for v := range ch {
		// pretend to persisting result
		fmt.Println(v)
	}
}

const MAX_THREAD = 3

func main() {
	inputChan := GetInputChan()

	// fan-out
	pipelines := make([]<-chan int, MAX_THREAD)
	for i := 0; i < MAX_THREAD; i++ {
		pipelines[i] = GeneratePipeline(inputChan)
	}

	// fan-in
	merged := MergePipelines(pipelines...)

	go PersistResult(merged)

	time.Sleep(10 * time.Second)
	fmt.Println("program exit")
}
```