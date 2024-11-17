# Go 有哪些 Reader 以及提供哪些读取操作？

## io

`Reader` 是 Go 中所有读操作的基础接口：

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

`Copy` 有缓冲，默认缓冲大小是 32KB；`ReadFull` 无缓冲，用于读取指定长度的字节。

## bufio

`Reader`: 结构体，默认的缓存大小 4KB。提供了很多实用的读取方法，例如 `ReadString` 可以自己指定分隔符, `ReadByte`, `ReadRune` 等。

`Scanner` 提供逐行读取功能。

## crypto/rand

`Reader` 不提供缓冲，因此如果要高频率进行小读取，[建议](https://tip.golang.org/doc/go1.19#:~:text=Read)用 `bufio.Reader` 包起来。
