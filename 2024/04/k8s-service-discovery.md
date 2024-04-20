# K8S 服务发现

## 现象

K8S 网络中的服务发现不可用 ，nslookup 显示：`server can’t find {SERVICE_NAME}.svc.cluster.local`。

## 排查及解决

服务发现（DNS）不可用，Service 域名无法解析 → 重启 k3s 无效 → 验证到 kube-system 命名空间下 deploy/coredns 的 53 端口是通的，证实并不是无法连接 dns 服务导致的故障 → 更新 coredns 镜像仍未解决，证明不是 dns 组件本身异常 → 到这里意识到排查方向可能出错了，dns 服务根本就没问题，有问题的是服务 → 定义一个新的服务和 service——可以访问！ → 比较 manifest 差异，最终定位到是 nodeSelector 导致 pg 端口不通的 → 原来调度到某个节点上的 Pod 都无法被访问 → 是节点故障了，而且是 nmcli network off 导致的 → nmcli network on 并重启机器解决。
