# zerobase 内存地址

## Q

- zerobase 地址是啥？
- 为什么 fmt.Printf 有可能导致变量的内存逃逸？

## A

逃逸到堆上的空值（`struct{}` 类型）都会指向同一个地址——Go 运行时的 zerobase 地址。

`fmt.Printf()` 可能会导致变量逃逸到堆上，因此可以写出以下测试代码：

```go
a := struct{}{}
b := struct{}{}
fmt.Printf("a's pointer: %p, b's pointer: %p, %v\n", &a, &b, &a == &b)
// a, b 可能发生逃逸，导致打印 true，0x1008a7eb8 即为 zerobase
// output: a's address: 0x1008a7eb8, b's address: 0x1008a7eb8, true
```

为什么 `fmt.Printf` 有可能导致变量的内存逃逸？——`fmt.Printf` 第二个参数是接口类型，内部使用了反射，因此是有可能产生逃逸。

这里另有份代码，可以看到，逃逸后空结构类型的值或者它的数组，都指向 zerobase：

```go
package main

import (
	"fmt"
	_ "unsafe"
)

//go:linkname zerobase runtime.zerobase
var zerobase uintptr

func main() {
	var s struct{}
	var a [42]struct{}

	fmt.Printf("zerobase = %p\n", &zerobase)
	fmt.Printf("       s = %p\n", &s)
	fmt.Printf("       a = %p\n", &a)
}
```