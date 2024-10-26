# 资源管理

## Pod 是如何被驱逐的

1. kubelet 读取 Cgroups 或者使用 cAdvisor 监控到的数据，计算 Eviction 阈值
2. 当节点上某项资源达到 Soft Eviction 阈值超过一段优雅时间，或者达到 Hard Eviction 阈值，kubelet 便会触发驱逐过程
3. 驱逐优先级取决于 Pod 的 QoS (Quality of Service) 类型，依次驱逐：`BestEffort` Pod → 超过了 `requests` 的 `Burstable` Pod → 超过了 `limits` 的 `Guaranteed` Pod

强烈建议将 DaemonSet 都设置为 Guaranteed QoS，因为这类 Pod 驱逐之后又要在相同 Pod 上重新创建，回收意义不大。

## 如何提升某 Pod 可用性

1. 使 QoS 为最高级别 Guaranteed，具体做法是使 Pod 中每一个容器的 `resources` 和 `limits` 相等，或者只设置 `limits` 不设置 `resources` 
2. 通过设置 `cpuset` 把容器绑定到某个或多个 CPU 核心，而不是像 `cpushare` 共享 CPU 的计算能力，减少操作系统的上下文切换，能显著提升性能。具体做法是在 Guaranteed QoS 的基础上，将 cpu 设置为整数值，比如 `"2"`，kubelet 就会给它分配 2 个独占的 CPU 核心
