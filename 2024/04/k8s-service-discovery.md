# K8S 服务发现

## 什么是 Endpoint？

被 Service selector 选中的 Pod 称为它的 Endpoints。但只包括处于 Running 状态且 readinessProbe 检查通过的 Pod。

也是 k8s 中的一种资源对象。一个 Service 对应一个 Endpoint 对象，对应多个 endpoint。

## 什么是 Headless Service？它与 ClusterIP Service 有何区别？

[这篇](https://stackoverflow.com/a/52713482)答案特别好。简单来说，ClusterIP Service 是为你提供一个 VIP，不关心具体访问哪个 Pod；而 Headless Service 则是让你通过 Pod 级别的域名访问，比如 `pod-1.mysvc.default.svc.cluster.local`，这些 DNS A 记录让我们可以稳定访问到具体的 Pod。StatefulSets 使用的就是 Headless Service。
