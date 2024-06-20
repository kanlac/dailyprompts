# iptables-legacy 导致的容器端口不可达

## 现象

麒麟系统下运行定制 docker 安装程序后，服务间访问可以 ping 通，但 nc 访问端口不通。

```
$ docker exec -it container_a nc -zv container_b 80
nc: container_b (10.242.0.3:80): Host is unreachable
```

尝试过：

- [x]  重启 docker
- [x]  重建网络
- [x]  重启 NetworkManager
- [x]  重启或关闭防火墙
- [x]  清空 iptables 规则并重启 docker

均无效，只有重启系统能够恢复正常。

## 期望

避免重启系统，容器间能够正常访问。

## 分析

- 执行 `iptables -L` 查看 iptables 规则，末尾有一行：
    
    ```
    # Warning: iptables-legacy tables present, use iptables-legacy to see them
    ```
    
- 定制 Docker 安装程序中包含了 iptables RPM 的更新

## 推断

iptables 更新后，iptables-legacy 规则还在生效，且无法通过常规方式查看和修改。

## 验证

初始化系统并查看 iptables 命令：

```
[root@localhost ~]# update-alternatives --display iptables
iptables - 状态为自动。
链接当前指向 /usr/sbin/iptables-legacy
/usr/sbin/iptables-legacy - priority 10
从 ip6tables：/usr/sbin/ip6tables-legacy
从 iptables-restore：/usr/sbin/iptables-legacy-restore
从 iptables-save：/usr/sbin/iptables-legacy-save
从 ip6tables-restore：/usr/sbin/ip6tables-legacy-restore
从 ip6tables-save：/usr/sbin/ip6tables-legacy-save
当前“最佳”版本是 /usr/sbin/iptables-legacy。
[root@localhost ~]# which iptables
/usr/sbin/iptables
[root@localhost ~]# ll /usr/sbin/iptables
lrwxrwxrwx 1 root root 26 10月  9  2023 /usr/sbin/iptables -> /etc/alternatives/iptables
[root@localhost ~]# ll /usr/sbin/iptables-legacy
lrwxrwxrwx 1 root root 20  3月  7  2021 /usr/sbin/iptables-legacy -> xtables-legacy-multi
```

或者更简单的，执行：

```bash
readlink -f $(command -v iptables)
```

输出：

```
/usr/sbin/xtables-legacy-multi
```

运行定制 Docker 安装程序（安装 iptables RPM）后，同样执行 `readlink` 命令，输出：

```
/usr/sbin/xtables-nft-multi
```

发现 iptables 命令指向的地址已变化，且 legacy 相关二进制已被移除，导致无法使用 iptables-legacy 查看或修改之前的规则。

使用包含 iptables-legacy 命令的容器可以查看遗留规则：

```bash
docker run --privileged -it --net host --rm m.daocloud.io/docker.io/wbitt/network-multitool bash
bogon:/# iptables -L
```

可以看到 Chain FORWARD 下有一条 `reject-with icmp-host-prohibited` 规则。删除之：

```bash
bogon:/# iptables-legacy -D FORWARD -j REJECT --reject-with icmp-host-prohibited
```

重新尝试访问容器端口，成功。

## 最终解决方案

```bash
local kylin_release=$(cat /etc/kylin-release 2>/dev/null)

if [[ "$kylin_release" == *" V10 "* ]]; then
    # 准备备份并清空 iptables 规则
    local backup_dir="/var/backups/iptables"
    local timestamp=$(date +%Y%m%d%H%M%S)
    local backup_file="$backup_dir/iptables-backup-$timestamp.rules"
    mkdir -p $backup_dir

    if [[ -n $(command -v iptables) ]]; then
        local iptables_real_path=$(readlink -f $(command -v iptables))

        if [[ $iptables_real_path == *"legacy"* ]]; then
            echo "iptables links to a legacy version: $iptables_real_path"
            echo "Backup iptables rules into: $backup_file"
            iptables-save > $backup_file

            echo "Set default strategy to ACCEPT"
            iptables -P INPUT ACCEPT
            iptables -P FORWARD ACCEPT
            iptables -P OUTPUT ACCEPT

            echo "Clean iptables rules"
            iptables -F
            iptables -t nat -F
            iptables -t mangle -F
            iptables -X
            iptables -t nat -X
            iptables -t mangle -X
        fi
    fi
fi
```
