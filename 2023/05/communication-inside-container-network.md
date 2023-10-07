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

## Fast demo：验证 Docker 网络内部的服务名访问

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
