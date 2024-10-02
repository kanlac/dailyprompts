# kube-proxy

## 主要工作

1. **网络转发**：kube-proxy 监听 Kubernetes API 中的 Service 和 Endpoint 资源的变化，当新的 Service 创建或更新时，它会动态更新系统的（iptables 或 IPVS）**网络规则**，这些网络规则控制 Service 的虚拟 IP 和真实的后端 Pod IP 之间的映射关系，确保服务通过虚拟 IP（即 Service 的 ClusterIP）能够与 Pod 通信。
2. **负载均衡**：kube-proxy 实现了基本的负载均衡功能，将流量均匀分布到 Service 关联的多个 Pod 上。

kube-proxy 的主要职责不包括 DNS 解析，DNS 解析在 Kubernetes 集群中是由 CoreDNS 或 kube-dns 服务来处理的。

## 代理模式

kube-proxy 在 Linux 节点下有多种代理模式。

2021 年已将性能较差的 Userspace 模式[淘汰](https://www.notion.so/kube-proxy-d1b1d8d1f43544a68bfac05e56b2e5c9?pvs=21)，该模式下，当请求到达 CusterIP，内核会转发给用户态的 kube-proxy 处理。

[现在](https://kubernetes.io/docs/reference/networking/virtual-ips/#proxy-modes)保留 3 种代理模式，分别用不同的方式配置包转发规则：

1. iptables - 对于 ClusterIP Service，默认写的规则是随机选择后端 Pod，也可以选择会话亲和性。NodePort 和 LoadBalancer 也差不多，但是会修改原始 IP
2. ipvs - IP 虚拟服务器，基于 netfilter hook function, 需要节点上支持 IPVS kernel modules。延迟更低，性能更好；提供更多负载均衡选项，包括 round robin，基于原始 IP 哈希等等
3. nftables - 基于内核的 netfilter subsystem，使用 nftables API 写规则，相当于下一代的 iptables API，性能好，属于比较新的技术
