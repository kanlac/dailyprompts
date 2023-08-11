# ENTRYPOINT v.s. CMD

### ENTRYPOINT 与 CMD 区别是什么？哪个先执行？
    
当 ENTRYPOINT 和 CMD 一起使用时，CMD 中的内容将作为参数传递给 ENTRYPOINT。例如，创建启动脚本 `entrypoint.sh`：

```bash
#!/bin/sh

# 这里是检查和修复 AOF 文件的代码

# 执行传递给脚本的任何其他命令
# 调研发现 redis 官方镜像用的也是这个
exec "$@"

```

和 Dockerfile：

```docker
ENTRYPOINT ["/entrypoint.sh"]
CMD ["redis-server"]

```

在上面的例子中，当容器启动时，会执行 `/entrypoint.sh redis-server`。如果在运行容器时覆盖了 CMD，例如 `docker run <image> redis-server --appendonly yes`，那么实际执行的命令将是 `/entrypoint.sh redis-server --appendonly yes`。

如果启动脚本没有 `exec "$@"` 这一行，CMD 就不会执行！

不过，如果没有写 ENTRYPOINT，那么 CMD 是会执行的！
    
### CMD 可以执行多个命令吗？
    
可以，使用 `sh -c` 是一个常见的技巧来在 `CMD` 或 `ENTRYPOINT` 中执行多个命令。你可以使用它将多个命令组合成一个单一的命令字符串。

例如：

```
FROM ubuntu:20.04
CMD sh -c "echo 'Hello from the first command' && echo 'Hello from the second command'"

```

这里，我们使用了 `sh -c` 以及 `&&` 来顺序执行两个 `echo` 命令。当你运行这个容器时，你会看到两个 `echo` 命令的输出。

请注意，使用 `sh -c` 要求您仔细处理命令字符串内的引号和特殊字符，以确保它们按预期的方式工作。