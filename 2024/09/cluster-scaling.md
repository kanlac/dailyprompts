# 集群伸缩

## Pod 伸缩

### 水平

控制 Pod 数量。

手动伸缩：kubectl scale, 或者 kubectl edit 修改 `spec.replicas` 字段的值。

自动伸缩：Horizontal Pod Autoscaler, HPA，比如可以自动让集群调整 Deployment Pod 数量：

1. 命令创建
    
    ```bash
    kubectl autoscale deployment example-deployment --cpu-percent=50 --min=2 --max=10
    ```
    
2. 定义清单
    
    ```yaml
    apiVersion: autoscaling/v2beta1
    kind: HorizontalPodAutoscaler
    metadata:
      name: example-hpa
    spec:
      scaleTargetRef:
        apiVersion: apps/v1
        kind: Deployment
        name: example-deployment
      minReplicas: 2
      maxReplicas: 10
      metrics:
      - type: Resource
        resource:
          name: cpu
          targetAverageUtilization: 50
    
    ```
    

### 垂直

调整资源用量。

1. 手动伸缩：`resources.requests` `resources.limits`
2. 自动伸缩：使用 Vertical Pod Autoscaler, VPA 可以自动调整 Pod 的 CPU 和内存请求。并不是原生支持的，需要单独安装 VPA 组件。示例 VPA 配置：
    
    ```yaml
    apiVersion: autoscaling.k8s.io/v1
    kind: VerticalPodAutoscaler
    metadata:
      name: example-vpa
    spec:
      targetRef:
        apiVersion: "apps/v1"
        kind: Deployment
        name: example-deployment
      updatePolicy:
        updateMode: "Auto"
    ```
    

## 节点伸缩

水平：节点数量；垂直：节点容量。

### 自动节点水平扩缩容 (Cluster Autoscaler)

自动调整集群中节点的数量。需要在云提供商的支持下配置。

## 实施扩容时的注意事项

考虑应用是否有状态，无状态可以直接扩容，有状态需要考虑扩容方案，比如使用 StatefulSet 等等

refer: https://kubernetes.io/docs/concepts/workloads/autoscaling/#scaling-workloads-vertically
