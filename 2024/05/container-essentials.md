# 容器核心技术

## 核心技术

Docker 使用 Linux 内核的一些特性实现它的基本功能：

### 1）Namespace 实现进程隔离

隔离的工作空间（isolated workspace）。当你创建一个容器时，Docker 会为它创建一系列的命名空间。

Docker Engine 使用的命名空间包括但不限于：

- `pid` 命名空间做进程隔离
- `net` 命名空间做网络隔离
- `ipc` 命名空间限制对进程间通信（Inter-Process Communication）资源的访问
- `mnt` 命名空间做文件系统挂载点管理

本质上是用 namespace 的 unshare 系统调用，告诉内核为子进程伪造虚拟的 PID 和网络等。

容器和虚拟机对进程的隔离有何不同？——容器中的进程对宿主机是可见的，只是 pid 不同。

创建一个容器执行了哪些关键的系统调用：

1. exec shim，shim 是 containerd 的子进程，容器内 init 进程的父进程
2. exec runc，runc 也就是 containerd 中实际运行容器的工具
3. unshare，实现资源隔离的关键，使子进程有独立的 namespace
4. clone 容器内 pid 1 进程，注意当我们从宿主机 strace，会看到这里返回的 pid 仍然是宿主机 pid，因为是在命名空间外调用的，但实际上它在自己的命名空间里是 pid 1
5. prctl 重命名为 init
6. 在容器内 clone 进程，这里会看到返回的 pid 是 2，是命名空间内的编号

`unshare` 系统调用是怎么做隔离的：

https://man7.org/linux/man-pages/man2/unshare.2.html 文档看 flag 介绍：

- CLONE_FS: Unshare filesystem attributes, so that the calling process no
longer shares its root directory (chroot(2)), current
directory (chdir(2)), or umask (umask(2)) attributes with
any other process.
- CLONE_NEWNET: Unshare the network namespace, so that the calling
process is moved into a new network namespace which is not
shared with any previously existing process. 实现网络隔离
- CLONE_NEWNS: unshare mount namespace, 实现存储隔离。mount namespace 管理进程可见的挂载点
- CLONE_NEWPID: Unshare the PID namespace, so that the calling
process has a new PID namespace for its children which is
not shared with any previously existing process.  The
calling process is not moved into the new namespace.  The
first child created by the calling process will have the
process ID 1 and will assume the role of init(1) in the
new namespace. 让子进程有自己的虚拟的 PID，容器内 init 进程 PID 为 1 的原因
- CLONE_NEWUSER: Unshare the user namespace, so that
the calling process is moved into a new user namespace
which is not shared with any previously existing process. 注意 docker 并没有用这个 flag，因此我们会看到容器内外的 uid 是一样的（只是用户名不同）

### 2）Cgroup 实现资源限制

cgroup（控制组）是 Linux 内核提供的一个功能，用于限制、监控和隔离进程组的资源使用。cgroup 可以对以下几种资源进行限制：

1. CPU，包括进程组的 CPU 时间，可以使用的核心等
2. 内存，包括限制进程组可以使用的最大内存，swap 空间等
3. 磁盘 I/O，包括限制进程组的读写速率
4. 网络，结合网络命名空间，可以限制进程组的网络带宽
5. 其他，比如设备，PID（进程组可以创建的最大进程数）等等

Cgroup 提供文件系统作为使用接口，在 /sys/fs/cgroup 下面编辑相应的文件，并填入 PID，Linux 就会将相关资源限制应用到进程。实际上服务资源限制（[资源管理](https://www.notion.so/730503cda44344848ac3da816264ecba?pvs=21) ）也是通过修改 Cgroup 文件系统来实现的。

### 3）联合文件系统（Union File System）

前面进程隔离已经提到过用 chroot 系统调用可以实现文件系统的隔离，不过 Docker 存储层使用的联合文件系统技术还是值得单独说一下。

分层镜像，也就是在 rootfs 的基础上，使用多个增量 rootfs 联合挂载一个完整的 rootfs 的方案，也就是多个文件系统层叠在一起，形成一个单一视图。解决了最大的依赖库——操作系统——的问题，实现了强一致性，打通了开发-测试-部署的每一个环节。

联合文件系统的优势：

1. 高效存储，多个镜像可以共享相同的基础层，从而减少重复数据的存储，节省磁盘空间
2. 一定程度上可以提高构建效率，只要不修改之前的层，就可以使用缓存
3. 快速启动，通过联合挂载可以迅速启动一个容器，不需要数据复制就能实现操作系统的「克隆」。容器运行时，Docker 会在最上层创建一个可写层，所有对文件系统的修改都会在这个层中进行，而不改变底层的只读层。参考《写时复制》in [并发问题](https://www.notion.so/6dce108777e14849a318e8ba11197d19?pvs=21) 。注意它是复制了操作系统但没有复制内核，所有容器的内核是共享的

有这么一个简单的小实验：把 bash 命令和它的依赖都拷贝到一个 test 目录下，并使用 `chroot` 命令（在 Docker 中则可能是 `pivot_root` ）修改 `bash` 二进制的根目录：

```go
chroot test /bin/bash
```

现在执行 `ls /` ，会看到它返回的都是 test 目录下的内容，说明它的根文件系统（rootfs）已经被修改。

## 容器与虚拟机的比较

容器相对虚拟机的优势是没有额外的性能损耗，缺点则是（通过 unshare namespace 来实现的）隔离不够彻底，用的还是相同的内核。比如在容器里执行 `top`，看到的依然是宿主机的内存和 CPU 的信息，这是因为 /proc 文件系统并不了解 Cgroup 的存在。又比如容器通过 `settimeofday(2)` 系统调用修改时间，那么整个宿主机的时间都会随之被修改。虽然有 Seccomp 等技术限制容器可以发起的系统调用，但是也不能完美地兼顾性能和隔离性。

