# K8s 监控体系

K8s 监控体系将 Prometheus 作为架构设计的核心。

## 三种指标来源

1. node exporter，通常作为 DaemonSet 运行，提供节点的负载、CPU、内存、磁盘和网络等信息
2. Kubelet 和 API Server 等组件的 /metrics，例如 API Server 的指标包含各个控制器的工作队列长度、请求的 QPS 和延迟数据等，这些都是检查 K8s 本身工作情况的重要依据
3. Metrics server，/apis，包含 Pod, Node, Container, Service 等核心监控数据，[不建议](https://groups.google.com/g/kubernetes-sig-instrumentation/c/bqdavcbwO7g/m/a9_4AU0DCAAJ)用于监控系统，因为它的指标主要用于自动扩缩容 pipeline，数据不太全面

容器的指标来自于集成到 kubelet 的 cAdvisor，通过 kubelet /metrics 或者 metrics server 都能获取到。

## Metrics Server 的部署机制

虽然看起来是通过 API Server 访问，但它独立部署，不是 API Server 的一部分。

它通过 kube-aggregator 插件的方式集成，我们访问 /apis 的时候，实际上访问的是一个代理，它后面有不同的指标后端。

需要[配置](https://kubernetes.io/docs/tasks/extend-kubernetes/configure-aggregation-layer/) API Server 开启 aggregator 功能才能正常使用。
