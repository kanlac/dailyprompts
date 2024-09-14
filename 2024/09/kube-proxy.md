# kube-proxy

## 主要工作

1. **网络转发**：kube-proxy 监听 Kubernetes API 中的 Service 和 Endpoint 资源的变化，当新的 Service 创建或更新时，它会动态更新系统的（iptables 或 IPVS）**网络规则**，这些网络规则控制 Service 的虚拟 IP 和真实的后端 Pod IP 之间的映射关系，确保服务通过虚拟 IP（即 Service 的 ClusterIP）能够与 Pod 通信。
2. **负载均衡**：kube-proxy 实现了基本的负载均衡功能，将流量均匀分布到 Service 关联的多个 Pod 上。

kube-proxy 的主要职责不包括 DNS 解析，DNS 解析在 Kubernetes 集群中是由 CoreDNS 或 kube-dns 服务来处理的。

## 如何实现负载均衡

kube-proxy 通过维护服务的虚拟 IP 到实际 Pod 的映射关系，并使用网络规则或内核模块实现负载均衡。有三种不同的代理模式，具体取决于集群的配置和操作系统支持的网络特性：

### (1) **Userspace 模式**（已过时，不推荐使用）

在 **Userspace 模式** 下，`kube-proxy` 运行在用户空间，拦截进入 Service 的请求，并手动选择一个后端 Pod，将请求转发给它。

- 工作流程：
    1. 当请求到达 Service 的虚拟 IP（ClusterIP）时，内核将请求交给 `kube-proxy` 的进程。
    2. `kube-proxy` 从该 Service 关联的后端 Pod 中选择一个 Pod，并将请求转发到该 Pod。
    3. `kube-proxy` 会均匀地选择后端 Pod，以实现基本的负载均衡。
- 缺点：
    1. 性能较低，因为所有的流量都需要通过 `kube-proxy` 的用户空间。
    2. 容易成为瓶颈，导致高延迟。

### (2) **iptables 模式**（当前的默认模式）

在 **iptables 模式** 下，`kube-proxy` 通过配置 Linux 内核中的 `iptables` 规则来实现负载均衡。

- 工作流程：
    1. `kube-proxy` 监控 Service 和 Pod 的变化，并为每个 Service 的 Endpoints 创建一系列 `iptables` 规则。
    2. 当流量到达某个 Service 的 ClusterIP 时，`iptables` 规则会将流量随机转发给该 Service 绑定的后端 Pod。
    3. 负载均衡是通过 `iptables` 的 **随机选择** 实现的，它将流量均匀地分发到不同的 Pod 上。
- 优点：
    1. **高性能**：所有流量直接由内核处理，不需要经过用户空间。
    2. **低延迟**：流量不经过 `kube-proxy`，减少了处理延迟。
    3. **高扩展性**：即使在大量请求下，负载均衡也能高效处理。

### (3) **IPVS 模式**（更高效的模式）

**IPVS（IP Virtual Server）模式** 是 `kube-proxy` 的最新代理模式，性能比 `iptables` 模式更好，特别适用于高流量场景。

- 工作流程：
    1. `kube-proxy` 通过监听 Service 和 Endpoints 的变化，使用 **IPVS** 内核模块创建路由表，将流量转发给后端的 Pod。
    2. 当请求到达 Service 的 ClusterIP 时，内核的 IPVS 模块根据配置的算法（如轮询、最小连接等）选择一个后端 Pod，转发流量。
- 支持的负载均衡算法：
    1. **轮询（rr）**：按顺序逐一将请求分发到不同的后端 Pod。
    2. **最小连接（lc）**：将流量转发到连接数最少的 Pod。
    3. **加权最小连接（wlc）**：考虑 Pod 权重，将流量转发到连接数较少且权重较大的 Pod。
    4. **加权轮询（wrr）**：按权重进行轮询。
- 优点：
    1. **高性能**：专门为负载均衡设计的内核模块，处理速度比 `iptables` 更快。
    2. **多种负载均衡算法**：可以根据不同的需求选择合适的调度算法。
    3. **动态性**：能更好地适应大规模的集群环境，快速响应后端 Pod 的变化。