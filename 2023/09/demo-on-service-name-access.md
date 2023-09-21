# 快速写一个 demo，验证 Docker 网络内部的服务名访问

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