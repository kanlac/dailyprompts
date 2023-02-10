# panic & recover

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

三、以下程序执行结果？
```go
func main() {
    defer func() {
      func() { fmt.Println(recover()) }()
    }()
    panic(1)
}
```

四、看程序回答问题
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

### 回答
一、会。当一个 goroutine 发生 panic，它会向它的父级 goroutine 抛出 panic，直到某个父级 goroutine 捕获（recover）。如果都没有 recover，则整个程序终止执行

二、`1, 2, panic`

三、直接 panic，不会输出 1，因为 recover 必须且只能套一层 func，不能多

四、B，因为 defer panic 会覆盖 panic