# 网络调试

### 如何关闭正在监听某一特定端口的进程
注意：不同操作系统不一样，比如 `netstat` 在 mac 上 `-p` 就不是 PID，而是 protocol

```bash
# 使用lsof
lsof -i :9093 
# 或者指定协议
lsof -i tcp:9093

# 使用 netstat
netstat -tuln | grep :8080

# （强制）结束进程
kill -9 1810
```

