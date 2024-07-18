# K8s Persistent Volumes

## PV

由管理员提供，它的特点：

1. 是存储的抽象，将存储的实现隐藏起来，可以是静态提供的 hostPath，也可以是动态提供的 NFS, iSCSI 或者云供应商特定的存储系统（通过存储类 Storage Class 实现）
2. （和 Volume 相比）有独立于 Pod 的生命周期

## PVC

由用户申请，用于将 PV 挂载到 Pod。类似 Pod 消费节点资源，PVC 消费的是 PV，Pod 可以请求特定级别的 CPU 或内存资源，PVC 可以请求特定大小和读写模式的 PV。

PVC 创建时，controlplane 会寻找匹配的 PV 进行绑定。每个 PVC 会绑定一个 PV，它们是一一对应的关系。匹配规则：

1. 如果 PVC 指定了存储类 Storage Class，则绑定的 PV 必须一致
2. PV 的容量必须大于等于 PVC——没规划好可能导致资源浪费
3. PV 的访问模式必须与 PVC 兼容

能否使用另一个命名空间的 PVC——否。

## **Storage Object in Use Protection**

当用户删除 PVC，它不会被立即删除，而是会推迟到没有 Pod 使用时再删除；当用户删除 PV，则会推迟到与 PVC 的绑定解除时删除。

删除正在使用的 PVC，会导致它进入 Terminating 状态，并能看到起保护作用的 Finalizer。

## StorageClass

1. `storageclass.kubernetes.io/is-default-class: "true"` 标注标记它为默认，如果有多个 StorageClass 被标记为默认，则 DefaultStorageClass 准入控制器会在创建 PVC 时报错
2. 绑定模式 `volumeBindingMode: WaitForFirstConsumer` 让 PVC 在 Pod 调度的时候创建，所以如果 Pod 通过 `nodeName` 绕过调度器会导致 PVC 一直 pending（如果通过标签指定节点亲和性就没问题）
3. `provisioner` 指定使用哪个云厂商的卷，若设置为 `kubernetes.io/no-provisioner` 则是本地卷

K3s 默认的 StorageClass：

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    defaultVolumeType: local
    ## ...
    storageclass.kubernetes.io/is-default-class: "true"
  name: local-path
provisioner: rancher.io/local-path
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer
```
