# kube-proxy

## 主要工作

1. **建立 Service 与 Pod 的映射**：kube-proxy 监听 Kubernetes API 中的 Service 和 Endpoint 资源的变化，当新的 Service 创建或更新时，它会动态更新系统的（iptables 或 IPVS）**网络规则**，这些网络规则控制 Service 的虚拟 IP 和真实的后端 Pod IP 之间的映射关系，确保服务通过虚拟 IP（即 Service 的 ClusterIP）能够与 Pod 通信。
2. **负载均衡**：kube-proxy 实现了基本的负载均衡功能，将流量均匀分布到 Service 关联的多个 Pod 上。

kube-proxy 的主要职责不包括 DNS 解析，DNS 解析在 Kubernetes 集群中是由 CoreDNS 或 kube-dns 服务来处理的。

## 代理模式

kube-proxy 在 Linux 节点下有多种代理模式。

2021 年已将性能较差的 Userspace 模式[淘汰](https://www.notion.so/kube-proxy-d1b1d8d1f43544a68bfac05e56b2e5c9?pvs=21)，该模式下，所有的规则写入和流量处理都是在用户态完成的。

[现在](https://kubernetes.io/docs/reference/networking/virtual-ips/#proxy-modes)保留 3 种代理模式，分别用不同的方式配置包转发规则：

1. iptables - 对于 ClusterIP Service，默认写的规则是随机选择后端 Pod，也可以选择会话亲和性。劣势：虽然流量处理是在内核态完成，但规则写入不是，当宿主机上有大量 Pod 时，kube-proxy 需要不断更新成百上千条 iptables 规则，会导致占用过多 CPU 资源
2. IPVS - IP 虚拟服务器，基于**内核模块**的网络流量处理技术，专注于进阶负载均衡。性能更好的部分原因：1）内核态写规则；2）使用更高效的数据结构（哈希表）。More：基于 netfilter hook function, 需要节点上支持 IPVS kernel modules。提供更多负载均衡选项，包括 round robin，基于原始 IP 哈希等等
3. nftables - 基于**内核模块**的网络流量处理技术，一种**网络过滤框架**，专注于数据包过滤和处理，是 iptables 的继任者。性能更好的部分原因：1）更高效的规则架构，减少了规则数目，提高了匹配速度；2）使用哈希表。More：netfilter subsystem，提供一个用户空间接口，允许用户使用 nftables API 写规则。注意 nftables 和 IPVS 是互补的技术不是替代

以 iptables 模式为例说明，当我们创建一个 Service，k8s 会为它分配一个虚拟 IP，这个 IP 只是用在 iptables 转发规则里，并没有一个实际的设备地址，因此 **ping 这个地址不会有任何响应**。

Service 转发时会不会修改原始 IP？——ClusterIP 不会，NodePort 和 LoadBalancer 会。
