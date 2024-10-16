# 发布工程 Release Engineering

## 从一些细节说明服务的规范化怎么做

配置管理：
- 配置应该集中管理，并且应该有一个清晰的版本控制系统。
- 敏感的配置信息（如密码和API密钥）应该加密，并且只在需要的时候解密。
- 配置的更改应该通过自动化的流程进行，以减少人为错误。

网络管理：
- 应该只暴露必要的端口，以减少安全风险。
- 所有的端口都应该有明确的用途，并且应该在文档中详细说明。
- 网络策略应该定期审查和更新，以应对新的安全威胁。

存储管理：
- 应该挂载必要的卷，以便于数据持久化和备份。
- 如果需要挂载额外的卷，考虑限制写入权限
- 数据的备份和恢复策略应该明确，并且应该定期测试。

容器管理：
- 在同一网络下，容器的名字应该是唯一的，以避免混淆。
- 容器的生命周期应该被正确管理，包括创建、更新、停止和删除。
- 容器的资源使用应该被监控，以便于发现和解决性能问题。

镜像管理：
- 尽量选择较小的公共镜像，以减少部署成本。
- 在构建镜像时，应该尽量减少镜像的大小，例如使用多阶段构建（multi-stage builds）。

## Google 的配置管理是怎么做的？

配置管理可能是不稳定性的一个重要来源，所以值得重视。

- 配置文件都放到代码仓库中，打包的时候使用主分支上的配置文件。开发者和 SRE 都可以修改这个文件
- 打包：一种方式是把配置文件和二进制文件打在同一个包里，这样做的好处是部署简单，只需要安装一个包；另一种方式是单独打一个配置包，和代码一样，用同一个系统编译和发布，每次构建有一个构建 ID，可以很方便地获取指定版本的配置。这样做的好处是可以分开管理，比如如果只需要修改一项配置，就只用重新构建并发布配置包，不需要重新安装二进制包。适合配置比较多比较复杂的场景。

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
