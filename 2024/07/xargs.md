# xargs

## 找出文件名中包含 "foo" 的文件并删除

```bash
# mac
find . -name 'foo' -print0 | xargs -0 rm -v
# linux
find . -name '*foo*' -print0 | xargs -0 rm -v
```

使用 `-print0` 和 `-0` 的原因是为了处理包含空格、换行符或其他特殊字符的文件名，将空字符（null character, `\0`）作为分隔符，而不是默认的换行符。

## 找出所有 redis 容器并打印它们的镜像哈希

```bash
docker ps --filter "name=redis" --format "{{.Names}}" | xargs -I {} docker inspect --format "{{.Name}}: {{.Image}}" {}
# 由于只需要在后续命令的末尾位置插入一个参数，所以其实可以省略 `-I`:
docker ps --filter "name=redis" --format "{{.Names}}" | xargs docker inspect --format "{{.Name}}: {{.Image}}"
```

`-I {}`：`-I` 选项指定一个替换字符串 `{}`。`xargs` 会将输入中的每一项用替换字符串标记，然后执行后面的命令。在这个例子中，`{}` 是替换字符串，它也可以是任何其他字符串，如 `%` 或 `@@`，但 `{}` 是最常见和直观的选择。
