# 容器资源隔离

## 创建一个容器执行了哪些系统调用

在宿主机上做 strace，可以找到以下关键系统调用：

1. exec shim，shim 是 containerd 的子进程，容器内 init 进程的父进程
2. exec runc，runc 也就是 containerd 中实际运行容器的工具
3. unshare，实现资源隔离的关键，使子进程有独立的 namespace
4. clone 容器内 pid 1 进程，注意当我们从宿主机 strace，会看到这里返回的 pid 仍然是宿主机 pid，因为是在命名空间外调用的，但实际上它在自己的命名空间里是 pid 1
5. prctl 重命名为 init
6. 在容器内 clone 进程，这里会看到返回的 pid 是 2，是命名空间内的编号

容器隔离本质上是用 namespace 的 unshare 系统调用，告诉内核为子进程组创建相应的命名空间。

## `unshare` 系统调用是怎么做隔离的

https://man7.org/linux/man-pages/man2/unshare.2.html 文档看 flag 介绍：

- CLONE_NEWNET: Unshare the network namespace, so that the calling
process is moved into a new network namespace which is not
shared with any previously existing process. 实现网络隔离
- CLONE_NEWPID: 让子进程有自己的虚拟的 PID，且为 1，这就是容器内 init 进程 PID 为 1 的原因
- CLONE_NEWIPC: 子进程组是否共享 IPC 资源。IPC：inter-process communication，进程间通信。Docker 会为容器创建单独的 IPC 命名空间，使它无法与其他外部进程通信
- CLONE_FS: Unshare filesystem attributes, so that the calling process no
longer shares its root directory (chroot(2)), current
directory (chdir(2)), or umask (umask(2)) attributes with
any other process.
- CLONE_NEWNS: unshare mount namespace, 实现存储隔离。mount namespace 管理进程可见的挂载点
- CLONE_NEWUSER: Unshare the user namespace, so that
the calling process is moved into a new user namespace
which is not shared with any previously existing process. 注意 docker 并没有用这个 flag，因此我们会看到容器内外的 uid 是一样的（只是用户名不同）

## 容器和虚拟机对进程的隔离有何不同

容器中的进程对宿主机是可见的，只是 pid 不同。
