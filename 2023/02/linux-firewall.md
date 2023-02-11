# Linux 防火墙相关问题

### 问题

一、如何查看指定 zone 下的防火墙规则？

二、firewalld 的配置文件位于什么地方？

三、如何开启防火墙端口？

### 回答

一、
```bash
firewall-cmd --zone=public --list-all
```

二、`/etc/firewalld`

三、
```bash
firewall-cmd --permanent --zone=public --add-port=25565/tcp --add-port=19132/udp
# 重启防火墙生效
firewall-cmd --reload
```