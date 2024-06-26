# 并发限制 (ulimit)

## Goroutine 的上限

没有严格限制，上限取决于机器性能，一般为千万级。但是，在实际的生产环境中，应该根据实际需求和系统资源来决定要启动多少 goroutine，以避免过度使用系统资源和导致系统崩溃。

## 文件描述符

### *nix 对文件描述符数量有怎样的限制？
    
通过 lsof 命令可以查看文件描述符。ulimit 限制的不是系统上总计文件描述符数量，也不是单个进程的，而是当前 shell 会话及其派生的子进程可开启的文件描述符的数量。
    
### 如何查看或修改限制
    
Linux操作系统中，每个进程对文件、套接字等资源的访问是通过文件描述符（file descriptor）来实现的。系统对每个进程所能使用的文件描述符数量有一定的限制。这些限制可以通过几种方式查询和修改。

查看限制的命令通常是`ulimit`，可以用于查看或设置用户级别的资源限制。

1. **软限制（Soft Limit）**：软限制是系统对资源使用的默认上限，可以被任意用户进程更改，只要新值不超过硬限制。
2. **硬限制（Hard Limit）**：硬限制是系统对资源使用的绝对上限，只有超级用户（root）可以更改。

查看文件打开数量限制：

```bash
# 查看硬限制
ulimit -Hn
# 查看软限制
ulimit -Sn
# 或
ulimit -n
```

默认限制通常设置为 256 或 1024。

如果你想要在当前shell会话中增加文件描述符的软限制，且有 root 权限，可以使用`ulimit -n <new limit>`命令。例如，如果你想把软限制设为2048，可以使用命令`ulimit -n 2048`。

然而，这只会改变当前shell会话的限制。如果你想要永久改变文件描述符的限制，你需要编辑系统级别的配置文件。

在大多数Linux系统中，系统级别的文件描述符限制存储在`/etc/security/limits.conf`文件中。你可以编辑这个文件来更改系统的文件描述符限制。此外，对于systemd管理的服务，还需要在相应的服务单元文件中设置`LimitNOFILE`。

需要注意的是，不同的Linux发行版可能有不同的方法来更改这个限制，因此你应该查看你的系统文档或使用搜索引擎来找到适用于你的系统的具体步骤。
    

## 用户进程

### 用户进程数量的上限
    
`ulimit -u` 命令用于查询或设置 Unix/Linux 系统中的用户进程的最大数量限制。`-u` 参数代表 "user processes"，即用户进程数。

ulimit 的默认值取决于系统的内存大小，一个示例值是 655350。在某些情况下，系统管理员可能会提高这个值以允许更多的并发进程。

你可以通过在命令行中执行 `ulimit -u` 命令查看当前的软限制。如果你需要修改这个值，你可能需要使用 `sudo` 权限，并使用 `ulimit -u [新的限制]` 命令来设定。注意，这种方式设置的新限制只在当前的 shell 会话中有效，如果需要永久修改，通常需要修改 `/etc/security/limits.conf` 文件或相应的系统配置文件。

如果你想要查看系统的所有`ulimit`值，你可以简单地运行`ulimit -a`。这将显示当前shell的所有资源限制。
    
### 如何用 Go 产生一个子进程？如何查看一个 Go 进程相关的所有进程？
    
在 Go 中，使用 `exec.Command` 会创建一个新的用户进程，也是它的子进程，可以通过 `ps -ef | grep <PID>` 看到。

```go
func main() {
    fmt.Println("Main process ID:", os.Getpid())

    // 启动一个子进程
    cmd := exec.Command("sleep", "60")
    if err := cmd.Start(); err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Child process ID:", cmd.Process.Pid)

    // 等待子进程完成
    time.Sleep(65 * time.Second)
}

```
    

## 单个 socket 或文件可以进行的并发操作数量上限

在大多数情况下，你可以在一个文件或套接字上进行多个并发读或写操作，但这并不意味着这些操作都会同时进行。相反，操作系统通常会对这些操作进行排队，并按照一定的顺序进行处理。

因此，一个文件或套接字可以“支持”的并发操作数量**实际上取决于操作系统如何处理并发I/O请求，以及它如何排队和调度这些请求**。这个数量并不固定，而是由操作系统的I/O子系统、硬件资源、文件系统和网络协议栈等多个因素共同决定的。

`panic: too many concurrent operations on a single file or socket (max 1048575)`

以上报错显示，Go 的并发操作上限为百万级别。

这个限制是**由 netpoller 的内部数据结构的设计所决定**的，该实现是Go的异步I/O机制的一部分。具体来说，如果你试图在单个文件或socket上同时执行超过1048575个操作，你会遇到你所提到的那个错误。

这个限制源于Go使用32位整数来追踪每个文件或socket上的并发操作。其中的一些位被用于其他目的，剩下的位用于追踪并发操作。在这种设置下，最大的并发操作数是1048575。

但是，需要注意的是，在实际应用中，大多数情况下不太可能达到这个限制，因为在单个文件或socket上同时执行这么多操作通常是不理想的。这种情况下，一个更好的解决方案可能是使用更多的文件或sockets，或者使用一种更有效的并发控制策略。
