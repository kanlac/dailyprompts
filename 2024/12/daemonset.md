# DaemonSet 的特性及其实现原理

[How Daemon Pods are scheduled](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/#how-daemon-pods-are-scheduled):

1. （Controller Manager 中的）DaemonSet 控制器定期监测集群中的节点状态，创建 Pod，设置**节点亲和性** `nodeAffinity`，加上一系列**容忍** `tolerations`，比如 `not-ready`，`unreachable` 等，如果 DS 使用主机网络（`spec.hostNetwork: true`），还会容忍 `network-unavailable` 污点，这适合用于部署网络插件
2. 调度器完成新 Pod 的绑定，设置 `nodeName`