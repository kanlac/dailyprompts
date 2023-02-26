# panic, recover and defer

### 问题
一、一个 goroutine 发生 panic 后，会影响其他 goroutine 执行吗？

二、以下程序执行结果？
```go
func main() {
  fmt.Print("1 ")
  defer recover()
  fmt.Print("2 ")
  var a []int
  _ = a[0]
  fmt.Print("3 ")
}
```

三、重复 panic
```go
func main() {
    defer func() { fmt.Println(recover()) }()
    defer func() { fmt.Println(recover()) }()
    defer panic(1)
    panic(2)
}
```
以上程序打印：
- A: 2 `<nil>`
- B: 1 `<nil>`
- C: 2 1
- D: 1 2
- E: 直接 panic

四、计算函数的运行时间
```go
func main() {
	startedAt := time.Now()
	defer fmt.Println(time.Since(startedAt))

	time.Sleep(time.Second)
}
```

五、defer 套娃
```go
func main() {
	defer fmt.Println("recover: ", recover())
	panic(1)
}
```

```go
func main() {
	defer func() {
		fmt.Println("recover: ", recover())
	}()
	panic(1)
}
```

```go
func main() {
	defer func() {
		func() { fmt.Println("recover: ", recover()) }()
	}()
	panic(1)
}
```


### 回答
一、会。当一个 goroutine 发生 panic，它会向它的父级 goroutine 抛出 panic，直到某个父级 goroutine 捕获（recover）。如果都没有 recover，则整个程序终止执行

二、`1, 2, panic`

三、B，因为 defer panic 会覆盖 panic

四、`125ns`，几乎为 0s。因为调用 defer 关键字会立刻拷贝函数中引用的外部参数，所以以 `time.Since(startedAt)` 的结果不是在 main 函数退出之前计算的，而是在 defer 关键字调用时计算的（[ref](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-defer/)）

五、
1. `<nil>`, `panic: 1`，捕获不到 panic
2. `recover: 1`
3. `<nil>`, `panic: 1`，捕获不到 panic