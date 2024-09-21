# Pod 调度

## 将 Pod 调度到特定节点

1. NodeName:
直接指定节点名称，这是最简单但也最不灵活的方式。
    
    ```yaml
    spec:
      nodeName: node1
    ```
    
2. NodeSelector:
使用节点标签来选择合适的节点。
    
    ```yaml
    spec:
      nodeSelector:
        disk: ssd
    ```
    
3. Affinity 和 Anti-Affinity:
提供更复杂和灵活的调度规则。有硬亲和（required）和软亲和（preferred）。
    - Node Affinity: 基于节点标签的更高级匹配。
    - Pod Affinity: 根据已经在运行的 Pod 来调度。
    - Pod Anti-Affinity: 避免与某些 Pod 在同一节点上运行。
    
    例如：
    
    ```yaml
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/e2e-az-name
                operator: In
                values:
                - e2e-az1
                - e2e-az2
    ```
    
4. Taints 和 Tolerations:
Taints 用于标记节点，而 Tolerations 允许 Pod 忽略这些 Taints。
    
    ```yaml
    spec:
      tolerations:
      - key: "key"
        operator: "Equal"
        value: "value"
        effect: "NoSchedule"
    ```
    
5. 自定义调度器:
开发自己的调度器来实现特定的调度逻辑。
6. Pod 拓扑分布约束:
控制 Pod 如何跨节点或集群拓扑分布。
    
    ```yaml
    spec:
      topologySpreadConstraints:
      - maxSkew: 1
        topologyKey: topology.kubernetes.io/zone
        whenUnsatisfiable: DoNotSchedule
        labelSelector:
          matchLabels:
            app: my-app
    ```
    
7. 污点和容忍度:
通过给节点添加污点并在 Pod 上定义相应的容忍度来控制调度。
8. 节点选择器策略:
在命名空间级别设置节点选择器要求。

## 如何均匀分布

### **Pod 亲和性和 `PodTopologySpreadConstraints`**

要在两个节点之间均匀地分配 5 个 Pod 的副本，可以使用 Kubernetes 的 `PodTopologySpreadConstraints`，它能指定 Pod 在多个节点之间尽量均匀分布。

### 示例 YAML 配置：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  replicas: 10
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      topologySpreadConstraints:
      - maxSkew: 1
        topologyKey: kubernetes.io/hostname
        whenUnsatisfiable: DoNotSchedule
        labelSelector:
          matchLabels:
            app: my-app
      containers:
      - name: my-container
        image: my-image
```

### 解释：

- **`topologySpreadConstraints`**: 用于指定 Pod 应在不同的拓扑单元（如节点、区域等）之间尽量均匀分布。
- **`maxSkew`**: 表示在每个节点上分布的 Pod 数量之间的最大差异。在这个例子中，`maxSkew: 1` 表示节点上最多可以有 1 个 Pod 数量的差异。
- **`topologyKey: kubernetes.io/hostname`**: 指定 Pod 应在不同的主机（节点）之间分布。
- **`whenUnsatisfiable: DoNotSchedule`**: 如果无法满足拓扑分布约束，Pod 将不会被调度。

通过这种方式，Kubernetes 将尽可能均匀地把 10 个 Pod 分配到两个节点上，每个节点上约有 5 个 Pod。