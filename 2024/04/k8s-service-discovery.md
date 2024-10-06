# K8S 服务发现

## 什么是 Endpoint？

被 Service selector 选中的 Pod 称为它的 Endpoints。但只包括处于 Running 状态且 readinessProbe 检查通过的 Pod。

## 什么是 Headless Service？它与 ClusterIP Service 有何区别？

[这篇](https://stackoverflow.com/a/52713482)答案特别好。简单来说，ClusterIP Service 是为你提供一个 VIP，不关心具体访问哪个 Pod；而Headless Service 则是（通过提供多个 DNS A 记录来）为 Pod 提供稳定的主机名。
