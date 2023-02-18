# Go 内存泄露场景

### Go 是否存在内存泄露问题？
Go 虽然理论上支持三色标记法自动管理内存，但是要理解，垃圾回收能够正常运行的前提是，程序中必须解除对内存的引用，这样垃圾回收才会将其判定为可回收内存并回收。

### 有哪些场景？
一）临时性泄露，指的是该释放的内存资源没有及时释放，对应的内存资源仍然有机会在更晚些时候被释放，即便如此在内存资源紧张情况下，也会是个问题。这类主要是 string、slice 底层 buffer 的错误共享，导致无用数据对象无法及时释放，或者 defer 函数导致的资源没有及时释放。

示例代码：
```go
func GetASlice() []int {
	s := make([]int, 10000000)  // 函数返回后，该空间不会被释放
	return s[:50]
}

func OperateOnSomeLongString(longString string) {
	s := longString[:5]
	fmt.Printf("do something with %s...", s)
	// 这里会导致临时性的内存泄露，因为 s 所引用的大字符串不会被回收
}
```

二）永久性泄露，指的是在进程后续生命周期内，泄露的内存都没有机会回收，如 goroutine 阻塞。Goroutine 会消耗内存和运行时资源，而 goroutine 的栈中引用着的堆数据不会被回收。Goroutine 本身是不会被 GC 的，它们必须自己退出，如果出现预期之外的 for-loop 或者 chan select-case 导致无法退出，就会导致协程栈及引用内存永久泄露问题。

### 如何避免临时性的内存泄露？（编写代码）

一）重置丢失的元素中的指针为 nil

二）使用 copy
```go
func CopyDigits(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    b = digitRegexp.Find(b)
    c := make([]byte, len(b))
    copy(c, b)
    return c
}
```
    
三）使用 append
```go
func CopyDigits(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    b = digitRegexp.Find(b)
    return append([]byte{}, b...)
}
```
---

参考：[Go程序内存泄露问题快速定位](https://www.hitzhangjie.pro/blog/2021-04-14-go程序内存泄露问题快速定位/)