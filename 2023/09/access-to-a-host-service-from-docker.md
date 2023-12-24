# 容器内如何访问宿主机端口

使用特殊的 DNS 域名 `host.docker.internal`。旧版本可能需要用 --add-host (docker run) 或者 extra_hosts (docker compose)。

旧版本 with docker compose：

```yaml
services:
  prometheus:
    image: bitnami/prometheus:2.48.0
    container_name: prometheus
    extra_hosts:
      - host.docker.internal:host-gateway
```

refer: [https://docs.docker.com/desktop/networking/](https://docs.docker.com/desktop/networking/)