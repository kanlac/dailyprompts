# seccomp

## 问题

在麒麟操作系统上 ZooKeeper 容器启动异常，发现以下日志：

```
[0.003s][warning][os,thread] Failed to start thread "GC Thread#0" - pthread_create failed (EPERM) for attributes: stacksize: 1024k, guardsize: 4k, detached.
```

- 分析
    
    **`EPERM`** 是系统调用一个错误代码，表示“Operation not permitted（操作不允许）”。这通常意味着进程试图执行一个它没有权限执行的操作。
    

## seccomp

`**seccomp`（secure computing mode）是 Linux 内核提供的一种安全机制，它可以限制进程可以调用的系统调用（system calls）**。当一个进程调用一个不被允许的系统调用时，`seccomp`可以配置不同的**“动作”（actions）**来决定如何处理这个调用。

`seccomp`的动作包括：

1. **SECCOMP_RET_KILL_PROCESS**：立即终止进程。
2. **SECCOMP_RET_KILL_THREAD**：立即终止调用线程。
3. **SECCOMP_RET_TRAP**：发送SIGSYS信号给调用的进程。
4. **SECCOMP_RET_ERRNO**：返回错误码，不执行系统调用。
5. **SECCOMP_RET_TRACE**：允许跟踪器决定是否允许或拒绝系统调用。
6. **SECCOMP_RET_ALLOW**：允许系统调用。
7. **SECCOMP_RET_LOG**：记录审计日志，然后允许系统调用。

这些动作允许开发者细粒度地控制进程可以执行的系统调用，增强系统的安全性。在创建安全容器或沙盒环境时，`seccomp`是一个非常有用的工具。

## 解决方案

docker compose 为相关服务添加以下配置，解决 JVM 进程问题

```yaml
security_opt:
      - seccomp:unconfined
```