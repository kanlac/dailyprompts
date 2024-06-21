# Debug: 容器无法停止/杀死

## 现象

重建容器失败：

```
tried to kill container, but did not receive an exit event.
```

无法停止容器，也无法删除，`docker kill -9` 也失败。

容器日志无异常。

## 排查

使用 strace 跟踪系统调用，捕获输出。中文输出可能需要禁止转换，比如从八进制转中文。

定位到是脚本中健康检查失败，执行了 kill 1。PID 为 1 的进程也就是 init 进程。

在宿主机上获取该容器的 init 进程。

```
ps aux | grep CONTAINER_HASH
```

`kill PID` 杀之。如果成功容器就重启了，但这里的问题是它无法被杀死。于是加上强制退出信号：`kill -9 PID`，容器终于停止。

也可以尝试杀死 shim 进程，也就是 init 进程的父进程：

```
ps -ef | grep INIT_PROCESS_PID
kill -9 SHIM_PROCESS_PID
```

## 解决

将健康检查脚本中的 `kill 1` 改为 `kill -9 1`。
