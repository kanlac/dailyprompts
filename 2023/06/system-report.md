# 查看系统报告

## sysstat 包含哪些工具？可以实现什么功能？

[sar 使用指南](sar-manual.md)

## 如何诊断导致服务器宕机(OOM)的具体进程？
- `dmesg` 查看内核环形缓冲队列，发现 OOM(out of memory) 错误
- `ps -p [PID] -o lstart` 对比进程的启动时间，定位到导致宕机的进程
- （让应用开发者）查看相对应的进程日志或容器日志，分析导致 OOM 的原因