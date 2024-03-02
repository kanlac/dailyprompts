# SAR 命令使用指南

sysstat 是一个用于查看和监控 Linux 系统性能的工具集，其中的 sar 是 System Activity Reporter 的缩写，它用于收集、报告和存储系统活动信息。

以下是一份简单的 SAR 命令使用指南：

## 检查 sar 是否启动

在许多 Linux 发行版中，sar 是作为 sysstat 包的一部分自动安装并启动的。您可以使用以下命令检查 sar 是否正在运行：

```bash
systemctl status sysstat
```

如果 `sar` 已启动，您应该会看到 `active (running)` 的状态信息。

## 基本使用

sar 命令格式如下：

```bash
sar [options] [interval] [times]
```

其中，`interval` 是采样间隔，`times` 是采样次数。

例如，如果你想每5秒钟采样一次，总共采样12次，可以使用如下命令：

```bash
sar 5 12
```

## 常见选项

以下是 SAR 命令中的一些常见选项：
- `-A`：报告所有可用的数据。
- `-u`：报告 CPU 使用率。
- `-r`：报告内存使用情况。
- `-b`：报告 I/O 活动。
- `-n DEV`：报告网络使用情况。

## 示例

```bash
# 默认显示当前整天的 CPU 使用情况
sar

# 查看所有可用的 SAR 数据
sar -A

# 每5秒钟报告一次 CPU 使用率，总共报告6次
sar -u 5 6
```

### 查看网络情况

```bash
sar -n DEV
```

您可以使用 `-n` 参数 搭配不同的关键词 (`DEV`, `EDEV`, `NFS`, `NFSD`, `SOCK`, `IP`, `EIP`, `ICMP`, `EICMP`, `TCP`, `ETCP`, `UDP`, `SOCK6`, `IP6`, `EIP6`, `ICMP6`, `EICMP6`, `UDP6`) 。

## 其他参数

`sar` 命令有很多其他参数，您可以查看 `man sar` 获取更多信息。

请注意，以上所有命令都需要 root 权限或 sudo 权限执行。

## 命令解释
下面是一个示例输出及其字段的解释：
```
Linux 4.18.0-25-generic (ubuntu)   06/02/2023  _x86_64_    (4 CPU)

12:10:01        CPU     %user     %nice   %system   %iowait    %steal     %idle
12:10:03        all     16.01      0.00      2.67      0.00      0.00     81.32
Average:        all     16.01      0.00      2.67      0.00      0.00     81.32
```

- %user: 用户空间中 CPU 时间的百分比，即应用程序使用的 CPU 时间的百分比。
- %nice: 以优先级调整为负数的用户空间 CPU 时间的百分比，即 CPU 时间被分配给高优先级任务的百分比。
- %system: 内核空间中 CPU 时间的百分比。
- %iowait: CPU 等待 I/O 操作完成的时间百分比。
- %steal: 虚拟环境中等待虚拟 CPU 资源的时间百分比。
- %idle: CPU 闲置的时间百分比。

## 查看特定时间段的 CPU 或 I/O 占用信息

`sar` 命令允许您指定要查看的日志文件，以便您可以查看特定时间段的历史数据。默认情况下，`sar` 命令将查看今天的数据，但是您可以通过使用 `-f` 参数来指定其他日期的日志文件。

例如，要查看前一天的数据，您可以运行以下命令：

```bash
sar -f /var/log/sa/sa$(date -d "yesterday" '+%d')
```

您也可以指定查看数据的时间范围。例如，要查看昨天 2 PM 到 4 PM 的数据，您可以运行以下命令：

```bash
sar -s 14:00:00 -e 16:00:00 -f /var/log/sa/sa$(date -d "yesterday" '+%d')
```

同样，可以使用 `-b` 参数来查看 I/O 数据：

```bash
sar -b -f /var/log/sa/sa$(date -d "yesterday" '+%d')
```

## 查找导致高占用的相关进程

当您确定系统在某个时间段内 CPU 或 I/O 使用率过高时，您可能希望找出导致高占用的相关进程。要做到这一点，您可以使用 `pidstat` 命令，这也是 `sysstat` 工具包的一部分。

`pidstat` 命令的基本语法如下：

```bash
pidstat [interval] [times]
```

例如，要查看每个进程每 5 秒的 CPU 使用情况，并报告 3 次，您可以运行以下命令：

```bash
pidstat -u 5 3
```

如果您想查看特定进程的 I/O 情况，您可以使用 `-p` 参数，后面跟上进程 ID。例如：

```bash
pidstat -d -p [PID] 5 3
```

## 查看特定进程在特定时间段的资源占用
虽然 `sar` 命令能够帮助我们分析特定时间段的系统资源使用情况，如 CPU、内存、I/O 等，但它并不能提供特定进程的资源使用信息。`sar` 收集的是系统级别的数据，主要用于追踪整个系统在一段时间内的性能变化。

如果你想查看特定进程在特定时间段的资源使用情况，你需要使用像 `pidstat` 这样的进程级别的性能监控工具。然而，这些工具主要提供实时的数据，而不是历史的数据。

为了追踪特定进程的资源使用情况，你可能需要定期运行这些命令并保存其输出结果。例如，你可以创建一个 cron 任务，每隔一定的时间运行一次 `pidstat`，并将结果保存到日志文件中。

```bash
# Crontab entry
* * * * * pidstat -u -p [PID] 1 59 >> /var/log/pidstat.log
```

以上 Cron 配置将每分钟记录一次特定进程的 CPU 使用情况，然后将这些信息追加到 `/var/log/pidstat.log` 文件中。然后，你就可以从这个日志文件中获取特定时间段内特定进程的性能数据。