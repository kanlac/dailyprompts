# Go 服务端取消请求

## HTTP 请求关闭时 Go 服务端的上下文是否会取消？有什么特例？

一般来说，关闭页面，刷新页面，或者 Ctrl+C 中止 curl 命令，都会关闭 TCP 连接，这样的话就 http 包的 context 就能正常 cancel。

但是有例外情况。

一）有请求转发，并且转发后的连接不会随着前面的连接关闭而关闭，这样就造成服务端感知不到连接关闭了。

二）没有读完 Body，或者没有 close body 的情况下，即使连接关闭也不会 cancel，这个比较坑，是 Go 实现的一个问题，17 年的 issue 一直悬而未决 https://github.com/golang/go/issues/23262