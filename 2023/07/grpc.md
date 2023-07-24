# gRPC

### gRPC 与 RESTful API 如何选择？（各自说出 3 个优势）
    
**使用 gRPC 的情况：**

1. 更高效的二进制数据序列化，比 JSON 包更小、解析速度更快
2. 需要双向流通信：如果你的应用需要服务器和客户端之间双向、实时的数据流，那么 gRPC 是一个不错的选择，因为它支持双向的流式 RPCs。
3. 支持 Unix 套接字：本地通讯更高效

**使用 RESTful API 的情况：**

1. 简单和易于使用：RESTful API 的设计和使用相对简单，并且广泛支持。它的学习曲线较低，且许多开发者都对其熟悉。
2. 浏览器兼容性：如果你的应用需要直接和浏览器交互，那么 RESTful API 可能是一个更好的选择，因为浏览器原生支持 HTTP 和 JSON。
3. 缓存需求：RESTful API 可以利用 HTTP 缓存机制来减少一些请求，改善性能。这在一些场景中（如获取静态、不经常变化的数据）可以很有用。
    
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