# Kubernetes Objects(Resources)

## ReplicaSet 的作用

Deployment 通过它来间接管理 Pod；每次更新 Deployment 会生成一个新的 ReplicaSet，旧的 RS 管理的 Pod 滚动删除，新的 RS 管理的 Pod 滚动创建，也就是说，RS 与 Deployment 的不同版本是一一对应的；通过 `kubectl describe` 命令可以看到相关 RS 的事件。

## DaemonSet 的特性及其实现

1. 在控制循环中确保每个节点一个 Pod，是通过**节点亲和性**实现的，创建 Pod 时会自动给 Pod 加上 `nodeAffinity`，确保这个 Pod 只会在这个节点上；
2. 会自动容忍有 `unschedulable` 污点的节点，这是通过**容忍性** `tolerations` 字段实现的。也可以配置容忍其他污点，比如 `network-unavailable` ，通过容忍这个污点，DS 会在尚未有可用容器网络的节点上面部署 Pod，这适用于部署网络插件。

## LoadBalancer Service 和 ClusterIP Service 有什么区别

后者只限于内部访问，而前者会生成 EXTERNAL-IP 可供外部访问，但它依赖于云提供商的负载均衡器实现

## LoadBalancer Service 和 HostPort Service 有什么区别

后者只是单纯地将服务在所在节点上开启端口，而前者做什么具体取决于负载均衡控制器的实现。对 K3s ServiceLB 来说，它也是在节点上开端口，不过还是有区别的：1）负载均衡控制器会给每个 service 创建一个对应的 pod，名称以 svclb 开头，在 kube-system 命名空间下；2）会经过负载均衡，访问某个节点的 IP 不代表访问该节点上的服务
