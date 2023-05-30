# Docker 有哪几种主要的容器网络模式？默认是哪种？
- bridge - 默认，容器间网络可以通过容器名访问
- host - 连接至宿主机的网络命名空间（可以用这个网络模式获取宿主机的 IP 地址）
- overlay - 在 Docker Swarm 集群中，Overlay 网络允许跨多个 Docker daemon 的容器相互通信
- macvlan - 连接至宿主机的物理网络，不常用

以下是默认的网络配置（等效于不配置）：
```yaml
networks:
  default:
    driver: bridge
```

如果您没有为服务明确指定网络，并且也没有定义名为 **`default`** 的网络，Docker Compose 将自动创建一个名为 **`default`** 的网络，并使用 **`bridge`** 驱动。