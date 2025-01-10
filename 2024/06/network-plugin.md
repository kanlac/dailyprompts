# 网络插件 Network Plugin

## 网络插件的作用

Flannel, Cilium 和 Calico 都是网络插件，它们的作用是把不同宿主机上的特殊设备（TUN 或 VETP）连通，在 k8s 中构建覆盖网络（overlay network），从而实现跨主机通信。

网络插件会为每个节点分配一个子网，可以通过`etcdctl` 命令查看宿主机与子网的关系，可以在节点上使用 `ip route` 命令查看路由交叉比较。

会给 Pod 分配 IP 地址。

## CNI 网络插件额外要做什么事情

用 cni0 网桥替代 docker0 网桥。两者功能完全一样。例如在部署 Flannel 的机器下，如果用 docker run，容器里的 eth0 设备会连接到主机上的 docker0 网桥（，docker0 再连接到 TUN 或 VETP），如果是 K8s 创建的就是连接到 cni0。

## K8s 如何配置 Pod 网络

CNI 的设计思想就是，K8s 在启动 Pod 中的 Infra 容器后，可以调用 CNI 网络插件为这个 Infra 容器的 Network Namespace 配置符合预期的网络栈。一个 Network Namespace 的网络栈包括网卡、回环设备、路由表和 iptables 规则。
