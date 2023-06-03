# IP 地址段配置/冲突解决

Q：docker-compose up 重启之后机器连不上了，可能是什么原因导致的？怎么解决？

A：使用 `netstat -rn` 查看路由表配置，看是不是网关与本地机器的网络冲突了。如果是，需要在 /etc/docker/daemon.json 中配置地址段。

k8s 在默认情况下会配置 Pod 和 Service 的 CIDR 网段，也可以在安装的时候通过选项指定。