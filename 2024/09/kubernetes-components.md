# Kubernetes Components

## 控制平面和工作节点的核心组件及功能

- Control Plane components
    - kube-apiserver｜核心，Control Plane 与各节点之间的网关
        - 具体使用途径
            - kubectl, kubeadm
            - REST 请求
            - client libraries for various programming languages, e.g. client-go
        - 采用 OpenAPI 提案
    - etcd｜持久化存储
        - 集群大脑
            - controller manager 通过 etcd 知道集群状态是否发生了变化
            - scheduler 通过 etcd 知道有哪些可用资源
            - 接收请求时，api server 通过 etcd 知道集群状态是否健康
        - 记录所有的集群操作，但不保存应用数据
    - kube-scheduler｜监控新创建的 pod，并决定将其分配给哪个节点运行（通过指示 kubelet）
    - kube-controller-manager｜检测节点状态变化，运行以下 controller 进程：
        - Node controller｜负责监控故障的节点
        - Job controller｜监视单次执行的任务，并创建 pod 以执行任务
        - Endpoints controller｜填充 Endpoints 对象
        - Service Account & Token controllers｜为新的命名空间创建新的账户和 token
    - cloud-controller-manager (optional)｜连接 Cloud Provider API，仅运行特定于云供应商的 controller（相对于集群控制器），包括
        - Node controller｜当节点失去响应后，它检查云供应商是否已将节点从云端删除
        - Route controller｜用于在底层云基础架构中设置路由
        - Service controller｜用于增删改云供应商的负载均衡
- Node components
    - kubelet｜启动（包含有容器运行时的）pod
        
        其工作原理是一个控制循环，而在这个 SyncLoop 上又包含很多子控制循环：
        
        1. 卷的管理
        2. Pod 生命周期的管理，包括调用 [CRI 容器运行时接口](https://www.notion.so/CRI-7c02be7fc015466193331ae73d69ab81?pvs=21)  接口创建容器
        3. 节点状态的管理，收集节点的状态信息上报给 API Server
        4. CPU 管理，维护该节点的 CPU 信息，以便 Pod 的 cpuset 功能能够使用
        5. ……
    - kube-proxy (optional)
    - Container Runtime｜支持 CRI 接口的运行容器的软件，如 Docker