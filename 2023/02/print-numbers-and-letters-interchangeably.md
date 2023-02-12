# 交替打印数字和字母

### 问题
使用 2 个 goroutine 打印数字和字母，最终输出要和以下一致：
> 12AB34CD56EF78GH910IJ1112KL1314MN1516OP1718QR1920ST2122UV2324WX2526YZ2728

### 解答
```go
package main

import (
	"fmt"
	"sync"
)

func printNumber(chNumber <-chan struct{}, chLetter chan<- struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	n := 1
	for range chNumber {
		fmt.Printf("%d", n)
		n++
		fmt.Printf("%d", n)
		n++

		chLetter <- struct{}{}
	}
	close(chLetter)
}

func printLetter(chLetter <-chan struct{}, chNumber chan<- struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	l := 'A'
	for range chLetter {
		if l >= 'Z' {
			close(chNumber)
			return
		}

		fmt.Printf("%c", l)
		l++
		fmt.Printf("%c", l)
		l++

		chNumber <- struct{}{}
	}
}

func main() {
	chNumber, chLetter := make(chan struct{}), make(chan struct{})
	var wg sync.WaitGroup
	go printNumber(chNumber, chLetter, &wg)
	go printLetter(chLetter, chNumber, &wg)
	chNumber <- struct{}{}
	wg.Wait()
}
```