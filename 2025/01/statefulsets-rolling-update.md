# StatefulSets 滚动更新

## StatefulSet 如何实现滚动更新？

使用 StatefulSet 的 `spec.updateStrategy.rollingUpdate.partition` 字段，比如 partition = 2，就是说所有序号大于等于 2 的 pod 会更新。可以在更新故障时维持系统的可用性，也可以用于将要被移除的依赖服务。

## StatefulSet 挂掉的 Pod 无法恢复是什么情况？怎么解决？

Answered on [StackOverflow](https://stackoverflow.com/a/79093393/8396777).

sts 默认的 Pod 管理策略是 `OrderedReady`，也就是会按照顺序，停掉了启动的 Pod 之后再启动新 Pod。但这会导致一个问题，就是如果 Pod 配置错误，比如镜像地址配置错误，导致进入 PullImageBackOff 的 Pending 状态，它就永远不会被替换，除非你手动删除这个 Pod。为了避免手动操作，可以选择放弃按顺序的 Pod 替换，也就是平行替换：

```yaml
apiVersion: apps/v1
kind: StatefulSet
spec:
  # Setting podManagementPolicy to Parallel is necessary; otherwise, the pod may become stuck and require manual intervention
  # refer: https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#forced-rollback
  podManagementPolicy: Parallel
```
