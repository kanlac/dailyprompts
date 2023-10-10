# seccomp

## 现象

在麒麟系统上安装应用，个别容器（例如 ZooKeeper）无法正常启动，查看日志发现 OOM（out of memory）错误。

进一步查看日志，看到 

```
[0.003s][warning][os,thread] Failed to start thread "GC Thread#0" - pthread_create failed (EPERM) for attributes: stacksize: 1024k, guardsize: 4k, detached.
```

## 分析

OOM 也好，GC 也好，这些都不是关键，关键是 `EPERM`，这是一个 Linux 内核错误，代表 operation not permitted，原因是没有系统调用的权限。

## seccomp

seccomp，或者安全计算（secure computing），是 Linux 内核的一个功能，用于限制系统调用。

Docker 内置一个[默认的 seccomp json 配置](https://github.com/moby/moby/blob/master/profiles/seccomp/default.json)，简单来说就是一个白名单。

## 解决方案

对于前述问题，完美的解决方案应该是基于 Docker 现有 seccomp 配置自行编写一份配置，不过能力有限，采用了折衷方案——在麒麟系统下关闭 seccomp。具体方法是在 /etc/docker/daemon.json 中添加：

```
"seccomp-profile": "unconfined"
```

## 快速检验

一个命令检测 EPERM 问题是否存在/已修复：

`docker run --rm --security-opt seccomp=seccomp-profile.json zookeeper:3.7.1-temuri`

## 参考

- docker seccomp 文档 https://docs.docker.com/engine/security/seccomp/
- docker 默认 seccomp 配置 https://github.com/moby/moby/blob/master/profiles/seccomp/default.json
- docker `security-opt` 选项文档 https://docs.docker.com/engine/reference/run/#security-configuration
- docker compose `security_opt` https://docs.docker.com/compose/compose-file/compose-file-v3/#security_opt
- docker compose 提供 seccomp 配置
    
    ```
    version: '3'
    services:
      web:
        image: nginx:latest
        security_opt:
          # 关闭 seccomp
          # - seccomp:unconfined
          # 或者，指定一份 seccomp 配置
          - seccomp:/path/to/seccomp/profile.json
    ```
