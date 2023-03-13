# 如何临时替换容器中的文件？如何修改镜像？

如何替换容器中的文件：

1. `docker cp` 命令可以在容器和宿主机之间拷贝文件
2. 使用 docker compose 的 `volumns` 特性挂载目录，配置好 yml 后 `docker compose up` 重启即可

如何修改镜像：`docker cp` 修改容器运行时后，可以通过 `docker commit` 命令可以保存新的镜像。