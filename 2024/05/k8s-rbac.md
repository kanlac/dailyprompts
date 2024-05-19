# K8S RBAC

k8s 里的 User 并不常用，更常用的扮演用户作用的对象是 ServiceAccount。创建一个 SA 时，k8s 会自动为其分配一个 Secret 对象。Secret 对象保存在 etcd 中，它包含用来跟 API Server 交互的授权文件，即 token。

## K8S 会为 Pod 挂载什么目录

Pod 示例：

```go
Name: zookeeper-0
Type: Pod
Containers:
  zookeeper:
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-qvd8l (ro)
```

我们可以在 Pod 的 spec 中定义它与哪个 ServiceAccount 绑定，如果不声明，k8s 会自动创建一个 default ServiceAccount 绑定给这个 Pod。描述 Pod 可以看到 SA 的 Secret 以目录的形式挂载到了 Pod 中的容器。思考：为什么每个 Pod 一定要自动自动挂载一个 Service Account Token？（即便它不需要与集群交互）——这是因为Kubernetes设计上倾向于提供一个统一的、开箱即用的环境，以适应多种可能的用例。这种默认行为确保了，如果Pod内的应用程序在未来需要与Kubernetes API进行交互，它们可以立即这么做，无需进行额外的配置。

## 如何找到某一 SA 有什么权限

需要找到包含这个 SA 的、命名空间下所有的 RoleBinding 和 ClusterRoleBinding，然后再查看对应的 Role/ClusterRole。

## ClusterRole(Binding) 与 Role(Binding) 的区别

是否跨命名空间，前者类别的资源是没有 namespace 字段的。
