# 网络调试

### 检查端口是否开放

`nc` 是 "netcat" 的简写，它是一个功能强大的网络工具，常被称为网络界的“瑞士军刀”。`nc` 可以用于构建各种网络连接，例如 TCP、UDP，并可以用于多种任务，例如端口扫描、文件传输、网络监听等。

有个很实用的技巧是用来检查容器 A 是否能（通过服务名）访问到容器 B 的某个端口：

```
docker exec -it CONTAINER_A sh -c "nc -zv CONTAINER_B 5432"
```

关于 `nc -zv` 的具体参数：

- `z`: Zero-I/O Mode，使用该选项让 `nc` 在扫描模式下运行，即它仅仅检查给定的端口是否可访问，而不是建立一个真正的连接。这在快速检查远程服务的可用性或本地系统上的开放端口时非常有用。
- `v`: 使 `nc` 运行在冗余模式，即输出更多的信息。在测试连接或端口扫描时，这可以帮助你看到连接尝试的具体结果。

### 关闭正在监听某一特定端口的进程

注意：不同操作系统不一样，比如 `netstat` 在 mac 上 `-p` 就不是 PID，而是 protocol

```bash
# 使用lsof
lsof -i :9093 
# 或者指定协议
lsof -i tcp:9093

# 使用 netstat
netstat -tulnp | grep :8080

# （强制）结束进程
kill -9 1810
```

- `lsof -i` 命令解释
    
    `lsof` 是 "list open files" 的缩写，也就是列出打开的文件。在 Unix 和 Linux 系统中，一切都是文件，包括网络套接字。
    
    - `i` 是 `lsof` 命令的一个选项，用于选择那些包含互联网地址的文件。
    
    当你使用 `lsof -i` 命令时，它会列出所有打开的网络连接，包括 TCP 和 UDP 连接。如果你在 `-i` 后面指定了一个端口号（例如，`lsof -i :8080`），那么它将只列出那些使用了指定端口的网络连接。
    
- `-tuln` 选项解释
    - `t` ：显示 TCP 连接。
    - `u` ：显示 UDP 连接。
    - `l` ：显示正在监听的套接字（即，显示服务器正在使用的网络服务和它们的状态）。
    - `n` ：以数字形式显示地址和端口号，不进行主机、用户名称以及网络服务名的解析。这可以加快命令的执行速度，因为它不需要查找这些名称。
    
    所以，`netstat -tuln` 命令会列出所有正在监听的 TCP 和 UDP 连接，并且以数字形式显示地址和端口号。
    
### 显示连接

`netstat -anp`

- `a` 查看所有协议的 socket
- `n` 将主机名显示为 IP，如果不加这个参数，获取主机名可能要很久
- `p` Linux 下表示显示 PID，Mac 下表示 protocol，需要额外参数

使用场景示例：

```bash
# 查看 state 为正在监听的 socket
netstat -an | grep LISTEN

# 查看 mysql 地址或端口
netstat -anp | grep mysql

# 查看某一端口是否被占用
netstat -tulnp | grep :80
# t: TCP, u: UDP, l: Listening, n: Numeric, p: process
```