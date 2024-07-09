# 四层端口注册服务

- 设计方案？包括如何实现动态增删，包括要给服务配置 cluster role K8S RBAC
- 原理？具体是怎么开放四层端口的？
- 困难？——数据隔离，基于节点 IP 做哈希。因为 SNAT 会导致丢失源 IP，使 service 配置的亲和性起不了实际作用，对比了 3 种方案但都不适用，包括 1）配置服务 Service 的亲和性；2）仅适用于外部均衡器的 `externalTrafficPolicy: Local`；3）仅适用于七层的 ingress 保留原始 IP 标注
- 解决方案？——通过 lua 实现基于节点 IP 做哈希
- 优化方向？——使用哈希环和虚拟节点实现一致性哈希 Consistent Hashing ，避免节点增删导致缓存丢失风暴；使用 bare-metal 四层负载均衡器，比如 MetalLB
