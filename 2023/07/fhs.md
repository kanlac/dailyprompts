# Filesystem Hierarchy Standard (FHS)

Q:
- /bin, /sbin, /usr/bin, /usr/local/bin 都有什么区别？
- 自己的脚本建议放在哪？
- docker, netstat, cp, scp, su, chown 这些命令分别位于哪个目录？

---

FHS 是 Linux 发行版的文件系统层级标准，由 Linux 基金会维护。

- `/bin` This directory contains executable programs which are needed in single user mode and to bring the system up or repair it. 包含单用户模式（一种超级用户模式）下需要用到的关键命令，包括启动和修复系统需要用到的 `cat`, `ls`, `cp` 等等
- `/sbin` Like `/bin`,  this  directory  holds commands needed to boot the system, but which are usually not executed by normal users. 包括 `fsck`, `init` 等
- `/usr` 适用于一般系统范围的二进制文件，只读不可写的用户数据，包含大部分的**多用户共享**的应用和工具
    - `/usr/bin` 普通用户可用，非关键性的命令二进制（单用户模式不需要），包括 `sed`, `scp`, `top`, `su` 等
    - `/usr/sbin` 系统管理员可用，但不需要用于启动或修复系统，包括各种网络服务的 daemon，包括 `lsof`, `chown`, `netstat` 等
    - `/usr/lib` 依赖库， `/usr/bin` 和 `/usr/sbin` 需要
    - `/usr/local` 该主机特定的，不受系统软件包管理，系统范围内可用
        - `/usr/local/bin` 自己的脚本通常建议放的位置，也包括 `docker`, `kubectl`, `redis` 等
        - `/usr/local/sbin`
        - `/usr/local/include` 安装 Protocol Buffers 时可以选择性地将一些知名的类型放到这个目录下
- `/boot` bootloader 文件
- `/dev` 设备文件，包括磁盘、键盘、摄像头等，比如一个硬盘是 `/dev/sda`，其中有分区 `/dev/sda1`, `/dev/sda2`…
- `/etc` 系统层级的配置文件，比如 apt, systemd, firewalld 应用的配置，还有 SSH 的系统级别配置
- `/home` 用户的目录
- `/lib` 依赖库， `/bin` 和 `/sbin` 需要用到
- `/media` CD-ROM 等可移除媒介的挂载点
- `/mnt` 临时挂载的文件系统，包括插上的U盘硬盘等
- `/opt` 可选（optional）的 add-on 的应用软件包，比如 Homebrew，供应商的应用也装在这里
- `/proc` 虚拟文件系统，将进程、CPU 和内核信息表示为文件，每个进程都有一个文件
- `/root` root 用户的 home 目录
- `/srv` web, FTP 等服务的目录
- `/tmp` 临时文件，比如有 Word 没有保存的文档，系统崩溃后就是从这里恢复，但通常重启后目录就清空了
- `/var` 变量文件，包含日志等