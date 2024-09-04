# Kubernetes Components

## 控制平面和工作节点的核心组件及功能

- Control Plane components
    - kube-apiserver｜核心，Control Plane 与各节点之间的网关
    - etcd｜持久化存储
    - kube-scheduler｜监控新创建的 pod，并决定将其分配给哪个节点运行（通过指示 kubelet）
    - kube-controller-manager｜检测节点状态变化，运行以下 controller 进程：
        - Node controller｜负责监控故障的节点
        - Job controller｜监视单次执行的任务，并创建 pod 以执行任务
        - Endpoints controller｜填充 Endpoints 对象
        - Service Account & Token controllers｜为新的命名空间创建新的账户和 token
    - cloud-controller-manager｜连接 Cloud Provider API，仅运行特定于云供应商的 controller（相对于集群控制器），包括
        - Node controller｜当节点失去响应后，它检查云供应商是否已将节点从云端删除
        - Route controller｜用于在底层云基础架构中设置路由
        - Service controller｜用于增删改云供应商的负载均衡
- Node components
    - kubelet｜启动（包含有容器运行时的）pod
    - kube proxy｜维护一些节点上的网络规则，以允许从集群外部或者内部与集群进行网络通信
        1. 是 Service 实现的关键组件，它监听 API Server 中 Service 和 Endpoints 的变化，并更新本地的网络规则，以允许从集群外部或者内部与集群进行网络通信
        2. 通过 CoreDNS（或其他 DNS 服务）完成 Service 名称的 DNS 解析
        3. Cilium 不仅提供完整的 CNI 实现，还可以替代 kube-proxy 的功能，使用 eBPF 技术更高效地处理网络包
    - Container Runtime｜支持 CRI 接口的运行容器的软件，如 Docker