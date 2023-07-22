# gRPC

### gRPC 与 RESTful API 有何不同？
    
是 Google 开源的 RPC 框架，使用 [Protocol Buffers](https://www.notion.so/Protocol-Buffers-0c833573286d440eb883fdf15d809831?pvs=21) （Protobuf）作为其接口定义语言（IDL），它支持多种编程语言，并提供了功能强大的双向流、流控、阻塞/非阻塞绑定等特性。
    
### gRPC 使用的协议是什么？为何选择这种协议？
    
gRPC 使用的是 HTTP/2 协议。HTTP/2 由于其二进制帧、头部压缩、多路复用和服务推送等特性，可以提供比 HTTP/1.x 更高的性能，减少延迟，提高网络利用率，这对于构建高性能的微服务非常有利。
    
### gRPC 可以定义哪几种类型的 service？
    
gRPC 提供了四种类型的服务：一元 RPC（Unary RPC, 类似于普通的函数调用）、服务端流式 RPC（客户端发送请求后，可以接收到服务端的多个响应）、客户端流式 RPC（客户端可以发送多个请求到服务端，然后等待服务端的响应）、双向流式 RPC（客户端和服务端都可以互相发送消息流，且各自的读写操作是相互独立的）。

服务端流式 RPC 定义如下：

```protobuf
// proto/sequence.proto
syntax = "proto3";

package sequence;

service Sequence {
    rpc Range (RangeRequest) returns (stream RangeResponse) {}
}

message RangeRequest {
    int32 start = 1;
    int32 end = 2;
}

message RangeResponse {
    int32 num = 1;
}
```
    
### gRPC 的拦截器有什么用？
    
gRPC 的拦截器（interceptor）是一种用于处理服务器和客户端之间调用的中间件，你可以使用它来处理或改变消息，例如验证、日志、限流、熔断、跟踪等。

例如，我们可以在 go 客户端拨号（`grpc.Dial`）时提供一个打印日志的拦截器。
    
### 在微服务架构中，使用 gRPC 有哪些优点和挑战？
    
在微服务架构中，使用 gRPC 有以下优点：**多语言支持**，**高效的二进制协议**，类型检查和编译时检查，**流控制**和错误处理等。挑战主要包括：debug 和跟踪的复杂性，防火墙和代理的兼容性问题，和学习曲线较陡等。