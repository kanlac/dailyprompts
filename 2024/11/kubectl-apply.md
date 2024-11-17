# kubectl apply

## Precaution before `k apply`?

在使用 `kubectl apply` 之前，最好检查一下现在集群和当前的 manifest 是否一致，这是为了避免集群使用 imperative commands 操作过，导致不再同步检查一致性的情况。

```
$ kubectl diff -f deployment.yaml
- replicas: 10
+ replicas: 5
```

不过因为 imperative command 不会留下记录，还是建议生产环境始终通过编辑 manifest 资源文件+搭配版本控制系统来控制集群。

## Byte length limitation

Why does `kubectl apply` has a byte length limitation while `kubectl create` doesn’t? How to solve it?

error:

```bash
Error from server (Invalid): error when creating "manifests/setup/0thanosrulerCustomResourceDefinition.yaml": CustomResourceDefinition.apiextensions.k8s.io "thanosrulers.monitoring.coreos.com" is invalid: metadata.annotations: Too long: must have at most 262144 byte
```

Solution: 

- `k apply --server-side`
- `k create something —dry-run=client -o yaml | k apply -f -`
