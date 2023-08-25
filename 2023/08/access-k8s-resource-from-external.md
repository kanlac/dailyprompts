# 从外部访问 k8s 资源有哪些方法

## **NodePort**

**`NodePort`** 是一种将集群内的服务暴露在每个节点上的特定端口上的方法。这样，任何能访问该节点的外部系统都能访问该服务。

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: my-app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: NodePort
```

## **LoadBalancer**

如果你在云提供商环境中运行 Kubernetes，你可以使用 **`LoadBalancer`** 类型的服务，该服务会自动创建云提供商的负载均衡器，并将流量路由到服务的 Pod。

这是在云环境（如 AWS, GCP, Azure 等）中公开服务的标准方法。

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: my-app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

## **Ingress**

使用 [Ingress](https://www.notion.so/Ingress-b47a62ebc9664969bfaa6e85da528d44?pvs=21) 控制器和资源，你可以更灵活地管理进入集群的流量。通过定义 Ingress 资源，你可以基于域名或 URL 路径来路由流量。

需要安装一个 Ingress 控制器，如 Nginx Ingress 控制器或 Traefik，它负责实现 Ingress 资源的规则；然后定义Ingress 资源。

## **HostNetwork**

通常，Pod 会在其自己的网络命名空间中运行，与主机网络隔离。但是，当设置 `hostNetwork=true` 时，Pod 会在主机的网络命名空间中运行，允许它监听主机的网络端口。

这种方法少用，因为它可能会造成网络配置上的冲突。

## **Port Forwarding**

使用 `kubectl port-forward` 命令，你可以将本地机器上的一个端口转发到集群内的 Pod 或服务。这主要用于开发或调试目的。

## **VPN / Direct Connect**

一些组织可能选择设置 VPN 或直接连接到 Kubernetes 集群的网络，以便安全地从外部访问。

### 总结

以上列出的每种方法都有其适用的场景和优缺点。选择哪种方法取决于你的需求，例如你是否在云上，你的安全需求，以及你希望为访问配置多少基础设施。