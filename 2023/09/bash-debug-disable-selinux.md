# A Bash Debug: Disable SELinux

## Q

SELinux 的配置文件如下：

```bash
ll /etc/sysconfig/ | grep selinux
# output:
# lrwxrwxrwx. 1 root root   19 Sep 11 21:45 selinux -> /etc/selinux/config
```

以下代码尝试持久化关闭 SELinux，它有什么问题？说明原因。

```bash
conf_path="/etc/sysconfig/selinux"
row_num=`grep -n "^SELINUX=" ${conf_path} | cut -d ":" -f 1`
sed -i -e "${row_num}c ${str}" ${conf_path}
```

## A

1. `sed -i` 会创建新文件，导致软链接被破坏，需要重建软链接
2. 应该修改软链接的目标文件，而不是修改软链接，否则强制 (force, `-f`)重建软链接后更改会被丢弃

修改后：

```bash
# edit target file
conf_path="/etc/selinux/config"
row_num=`grep -n "^SELINUX=" ${conf_path} | cut -d ":" -f 1`
sed -i -e "${row_num}c ${str}" ${conf_path}
# rebuild symbolic link after `sed -i`
ln -sf "${conf_path}" /etc/sysconfig/selinux
```