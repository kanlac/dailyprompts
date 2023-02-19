# 用流水线模式改写程序

[serial.go](https://go.dev/blog/pipelines/serial.go) 是一个 MD5 计算程序，它读取参数中的目录，输出目录下各文件的 MD5 值，并按路径名排序。没有任何并发操作。

```
% go run serial.go .
d47c2bbc28298ca9befdfbc5d3aa4e65  bounded.go
ee869afd31f83cbb2d10ee81b2b831dc  parallel.go
b88175e65fdcbc01ac08aaf1fd9b5e96  serial.go
```

现在，请运用流水线模式修改这段程序，第一阶段 `walkFiles` 遍历文件目录树，返回文件路径的 channel。

第二阶段 `digester` 计算文件的 MD5 摘要， 将 `result` 对象发送到 channel：

```
type result struct {
    path string
    sum  [md5.Size]byte
    err  error
}
```

最后一阶段，接收计算结果，并打印。

注意考虑错误处理，优雅退出和并发数限制。

---

答案参考：https://go.dev/blog/pipelines