# dockerd 访问权限

## 现象

```
permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock
```

## 查看 docker socket 属于哪个用户和用户组

```bash
ls -l /var/run/docker.sock
```

一般来说，centos 下默认属于 root 用户，docker 用户组。但实践发现在某些未知情况下，docker socket 属于 root 用户组。以下命令可以修改用户和用户组：

```bash
sudo chown root:docker /var/run/docker.sock
```

## 查看用户组里有哪些用户

```bash
# 查看 docker 用户组里有哪些用户
# 若 docker socket 属于 docker 用户组，
# 则以下命令列举出来的这些用户将具有访问 docker 的权限
getent group docker

# 将要执行 docker 命令的用户添加到 docker 用户组
# 使其具有访问 dockerd 的权限
sudo usermod -aG docker jack
```