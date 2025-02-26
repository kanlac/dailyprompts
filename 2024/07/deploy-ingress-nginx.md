# 如何部署 Ingress-Nginx

## 本地或云环境？

- 本地
    - 裸金属——主机网络或用 MetalLB 做负载均衡
    - 部署外部负载均衡器
- 云环境

## 是否接受单点故障？

- 是——使用主机网络，工作负载可以用 NodePort 或配置硬反亲和性的 Deployment
- 否——使用 MetalLB 提供虚拟外部 IP 作为高可用访问入口。MetalLB 可在二层工作，需要预留专用的 IP 地址池，不能用节点 IP

## 选择何种 Kubernetes Service？

- LoadBalancer——云环境或 MetalLB 支持的裸金属环境
- NodePort——本地开发或需要自定义外部负载均衡器
    - 是否固定端口？
        - 是——方便直接通过端口访问
        - 否——自动分配端口避免冲突
- 不使用（主机网络）——本地开发，或有极致性能需求（绕过 kube-proxy）