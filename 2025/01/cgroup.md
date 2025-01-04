# Cgroup 资源限制

cgroup（控制组）是 Linux 内核提供的一个功能，用于限制、监控和隔离进程组的资源使用。

cgroup 可以限制哪些资源：

1. CPU，包括进程组的 CPU 时间，可以使用的核心等
2. 内存，包括限制进程组可以使用的最大内存，swap 空间等
3. 磁盘 I/O，docker run 有[选项](https://docs.docker.com/engine/containers/run/#runtime-constraints-on-resources)可以限制 IO 读写速率
4. 其他，比如设备，PID（进程组可以创建的最大进程数）等等

注意 Cgroup 本身没法做到限制网络下载速率，但有[其他方法](https://www.reddit.com/r/docker/comments/1244rm7/is_there_an_easy_way_to_limit_download_speed_of_a/)可以做到这一点，比如可以在主机上使用 `tc`（Traffic Control）命令来限制容器的网络带宽。需要在容器的网络接口上添加带宽限制规则。

Cgroup 提供文件系统作为使用接口，在 /sys/fs/cgroup 下面编辑相应的文件，并填入 PID，Linux 就会将相关资源限制应用到进程。实际上服务资源限制也是通过修改 Cgroup 文件系统来实现的。
