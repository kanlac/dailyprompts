# Ingress-Nginx On Bare Metal K8s

## 方案一：MetalLB

给 ingress-nginx 配置 MetalLB，要[在二层模式下使用 MetalLB](https://metallb.universe.tf/concepts/layer2/)，需要预留专用的 IP 地址池，不能用节点 IP。ingress-nginx Service 使用 LoadBalancer 类型，MetalLB 会为它分配一个外部 IP，这个 IP 指向集群中的某一个节点，再经过它负载均衡到其他节点。

亮点是唯一入口，但尚不确定是否可以实现静态 IP。

为了保留源 IP，需要设置 traffic policy 为 local。

## 方案二：NodePort

使用 NodePort Service，客户端使用节点 IP 直接访问到节点上的容器，在进入 ingress-nginx 之前不会经过负载均衡。这种方式最简单。

注意如果访问的节点上没有 nginx 实例，包会被丢弃。

虽然实际上可以指定 NodePort 目标端口使用 80/443，但这样做意义不大，因为使用 NodePort 主要是为了自动分配端口避免冲突，如果固定要用 80/443，不如用主机网络。

## 方案三：主机网络

使 ingress Pod 跑在主机网络下，不需要 Service。为了避免自己端口冲突，同样需要用 DaemonSet。

## 方案四：externalIPs

ingress-nginx Pod 在容器网络内监听 80/443，但使用 K8s Service 的 externalIPs 功能可以将主机上 80/443 端口的流量转发到 Pod。唯独这种方式无法保留原始 IP，因此不建议使用。

---

refer: https://kubernetes.github.io/ingress-nginx/deploy/baremetal/#over-a-nodeport-service
