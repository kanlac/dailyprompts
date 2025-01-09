# 容器跨主机网络

## Flannel 的作用

给 Pod 分配 IP 地址，并在 k8s 中构建覆盖网络（overlay network），实现容器跨主机网络通信。

Flannel 为每个节点分配一个子网，可以通过`etcdctl` 命令查看宿主机与子网的关系，可以在节点上使用 `ip route` 命令查看路由交叉比较。

每个宿主机上的 flanneld 进程都监听 8252 端口。

## Flannel 有哪些后端实现

1. VXLAN, virtual extensible LAN
2. host-gw
3. UDP，因为性能问题已经被废弃

因为 UDP 模式最直接最易理解，所以书里从这个开始讲。

## VXLAN 技术相比 UDP 有什么优势

UDP 方案下，Node1 上的容器访问 Node2 上的容器，经过如下过程：

1. 容器出来的网络包，从用户态到内核态（第一次切换），先尝试给 docker0 处理，但因为不在 docker0 网段，所以会从 docker0 出来到宿主机上
2. 匹配到了 flannel0 设备，它是一个 TUN 设备，交由 flanneld 进程处理，从内核态到了用户态（第二次切换）
3. flanneld 在 etcd 中查找到目标主机的地址，封装到 UDP 包中经由 eth0 发送到 Node2 的 8252 端口（第三次切换）

Node2 上 flanneld 进程从 UDP 包中解封出 IP 包，发送到 flannel0 设备，后面的过程跟发出时正好相反。

从上面过程可以看到，光是将 IP 包发出就需要经过三次内核态用户态之间的切换，再加上 UDP 包的封装和解封都是在用户态完成，因此说 UDP 模式有严重的性能问题。

VXLAN 方案下的转发过程：

1. 容器出来的网络包，从用户态到内核态（第一次切换），先尝试给 docker0 处理，但因为不在 docker0 网段，所以会从 docker0 出来到宿主机上
2. 匹配到 flannel.1（VTEP 设备），在二层完成 UDP 封包并发出

主流的 VXLAN 为什么高效？因为发出只进行了一次上下文切换，且把包的封装和解封都放到内核态完成。

## TUN 设备 v.s. VTEP 设备

TUN 设备是一种**在内核态和用户态之间传递 IP 包**的网络设备，它在**三层（网络层）**工作。

VTEP 设备，VXLAN tunnel end point，在**二层（数据层 Ethernet）**工作，可以在**内核态**完成 UDP 的封包和解封，跟 UDP 模式下 flanneld 进程的作用非常类似。
