# IP 地址段配置/冲突解决

## 一
docker-compose up 重启之后机器连不上了，可能是什么原因导致的？怎么解决？

使用 `netstat -rn` 查看路由表配置，看是不是网关与本地机器的网络冲突了。如果是，需要在 /etc/docker/daemon.json 中配置地址段。

k8s 在默认情况下会配置 Pod 和 Service 的 CIDR 网段，也可以在安装的时候通过选项指定。

## 二
解释以下 docker 配置

```
{
    "bip": "10.241.185.0/23",
    "fixed-cidr": "10.241.185.0/24",
    "default-address-pools": [
        {
            "base": "10.240.0.0/12",
            "size": 16
        },
        {
            "base": "10.232.0.0/13",
            "size": 16
        },
        {
            "base": "10.230.0.0/15",
            "size": 16
        }
    ],
    "live-restore": true
}
```

**一句话解释：bip 是 docker 默认桥接网络的地址段，fixed-cider 是显示 bip 中实际给容器分配的地址段，default-address-pools 是新建 docker 网络的地址段**

冗长解释：

- `"bip": "10.241.185.0/23"`: 这个参数设置了默认的 Docker 网桥 `docker0` 的 IP 地址和子网。在这个例子中，`docker0` 将会使用 `10.241.185.0/23` 的网段。
- `"fixed-cidr": "10.241.185.0/24"`: 这个参数用 ‣ 设置了 `docker0` 网桥上可以为容器分配的 IP 地址范围。在这个例子中，Docker 将只会从 `10.241.185.0/24` 这个子网范围内分配 IP 地址给容器。注意这个网段必须是 `bip` 参数设置的网段的子集。
- `"default-address-pools"`: 这个参数定义了一组默认的地址池，这些地址池用于 Docker 为自定义网络分配子网。在这个例子中，有三个地址池，分别以 `10.240.0.0/12`，`10.232.0.0/13` 和 `10.230.0.0/15` 为基础，每个池的大小都是 16。这意味着，当您使用 `docker network create` 命令创建新的网络，但没有指定 `--subnet` 参数时，Docker 将会从这些地址池中选择一个子网给新的网络使用。
- `"live-restore": true`: 这个参数使 Docker 在重启时能够保留运行中的容器。这意味着，即使 Docker daemon 重启，您的容器仍然会继续运行，而且在 Docker daemon 重启后仍然可以管理这些容器。这个功能在需要更新或维护 Docker daemon 时非常有用，因为它可以减少对正在运行的容器的影响。

总的来说，这个 Docker 配置文件设定了 Docker daemon 的一些网络和运行行为参数，使 Docker 的运行能够满足特定的需求。

## 三

配置 /etc/docker/daemon.json 中的 default-address-pools 后，是否会影响 docker create network 命令创建的网关的路由？之前创建的网络是否也会受影响？如果不会，应该如何操作使其也在配置的网段范围内？

A：Docker 的 `default-address-pools` 选项允许你为 Docker 自动创建的网络定义一个默认的地址池。当你创建一个新的网络，但没有明确指定子网（`--subnet`）时，Docker 将从 `default-address-pools` 定义的范围中为新网络分配地址。

注意，`default-address-pools` 只对创建**之后**的网络有效，对创建**之前**的网络不会有任何影响。已经创建的网络已经被分配了特定的子网和网关，这些设置在网络创建时就已经确定，无法再修改。如果你希望之前创建的网络也使用新的地址池，你需要删除并重新创建这些网络。

例如，假设你在 `/etc/docker/daemon.json` 文件中配置了如下的 `default-address-pools`：

```json
{
  "default-address-pools": [
    {"base": "172.80.0.0/16", "size": 24}
  ]
}

```

然后你运行 `docker network create my_network`，Docker 将从 172.80.0.0/16 这个地址池中分配一个 /24 的子网给 `my_network`。如果你创建更多的网络，Docker 将继续从这个地址池中分配新的子网。

如果你希望修改已经创建的网络的子网，你需要先删除网络，然后再重新创建：

```bash
docker network rm my_old_network
docker network create --subnet=172.80.1.0/24 my_old_network
```

在这个例子中，`my_old_network` 将被分配一个新的子网 172.80.1.0/24。注意，当你删除一个网络时，所有连接到这个网络的容器都将失去网络连接，你可能需要停止这些容器，然后在新网络创建后再重新启动。