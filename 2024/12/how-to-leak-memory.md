# 有哪几种导致内存泄露的方式

## 示例 1：未清理的 Goroutine

```go
package main

import (
	"fmt"
	"time"
)

func leakGoroutine() {
	for {
		time.Sleep(1 * time.Second)
	}
}

func main() {
	for i := 0; i < 10; i++ {
		go leakGoroutine() // 启动多个未结束的 goroutine
	}
	time.Sleep(10 * time.Second) // 主程序睡眠，导致泄漏
	fmt.Println("Done")
}

```

## 示例 2：未释放的切片引用

```go
package main

import "fmt"

func main() {
	slice := make([]*int, 0)

	for i := 0; i < 100000; i++ {
		num := i
		slice = append(slice, &num) // 引用 num，造成泄漏
	}

	fmt.Println("Slice length:", len(slice))
}

```

## 示例 3：全局变量未清理

```go
package main

import "fmt"

var cache = make(map[string]string)

func storeValue(key, value string) {
	cache[key] = value // 全局变量 cache 不会被清理
}

func main() {
	for i := 0; i < 100000; i++ {
		storeValue(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	fmt.Println("Stored values in cache:", len(cache))
}

```

## 示例 4：循环引用

```go
package main

import "fmt"

type Node struct {
	value string
	next  *Node
}

func main() {
	node1 := &Node{value: "first"}
	node2 := &Node{value: "second"}

	node1.next = node2
	node2.next = node1 // 造成循环引用

	// 没有清理的情况下，node1 和 node2 无法被垃圾回收
	fmt.Println("Nodes created")
}

```