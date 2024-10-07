# CRI 容器运行时接口

## dockershim 的作用

每个节点上都会有一个负责响应 CRI gRPC 接口的 shim（垫片），比如 docker 就要实现一个 dockershim，它的作用就是把 CRI 请求翻译成具体的容器运行时的请求。

## crictl 为什么能查看 Pod

CRI 接口分为容器运行时接口和镜像接口，其中前者包含一个 RunPodSandbox 接口。为什么容器运行时需要关心 Pod 呢？这是因为当一个 Pod 里有多个容器时，不同的容器运行时的实现方式不同。比如 docker 会创建一个额外的 infra 容器，用来 hold 整个 Pod 的网络命名空间（参考 pause 容器），而 Kata Container 则会直接创建一个轻量级虚拟机充当 Pod。

## kubectl exec 的工作原理

首先将请求交给 API Server，其调用 kubelet 的 Exec API，kubelet 再从 shim 拿到一个 URL（shim 的 streaming server 的地址和端口），把它以 redirect 的方式返给 API Server，最后 API Server 通过重定向来向 streaming server 发起真正的 exec 请求，与它建立长连接。
