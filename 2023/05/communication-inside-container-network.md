# 容器网络内部的通信

## 容器之间的名字域名解析是怎么做的？

在同一个 Docker 网络中，容器之间可通过容器名通信。Docker 内置的域名解析服务器地址是 127.0.0.11。

如果使用 Nginx，需要配置这个 NS，否则容器间无法通过容器名通信：

```
location ~ ^/linglong/?(.*)$ {
    resolver 127.0.0.11 ipv6=off valid=3s;
		proxy_pass http://linglong_frontend/$1$is_args$args;
		# ...
}
```

## 如何利用现有的容器快速验证服务名访问是否可用？

例：`docker exec -it SOME_CONTAINER /bin/sh -c "nc -zv postgres 5432"`。`nc` 是 busybox 中包含的工具。

## 写一个 demo 验证服务名访问可用

虽然 chatgpt 可以给我弄出来，但是掌握这些基础也是很有必要的

```yaml
version: '3.8'

services:
  web:
    image: nginx:alpine

  client:
    image: busybox
    command: /bin/sh -c "wget -O- http://web"
    depends_on:
      - web
```

docker-compose up 启动，能看到 nginx 返回 html 表示访问成功。

备注：

- `wget http://web`: 下载到内容到当前目录
- `wget -O download http://web`: 保存到指定文件
- `wget -O- http://web`: 下载内容输出到标准输出
- 使用 `nc` 命令可以检测端口是否开启

## 服务名访问不可用时如何解决？

1. 系统内核参数 net.ipv4.ip_forward 确保为 1，修改后可能需要重启机器才能生效
2. 修改 docker daemon 配置
    
    在 /lib/systemd/system/docker.service 文件中添加：
    
    ```
    ExecStartPost=/usr/sbin/iptables -P FORWARD ACCEPT
    ```
    