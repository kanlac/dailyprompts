# 如何执行目录下的文本搜索，同时打印附近的行与行号？

```bash
grep -rnH -C 2 "foo" .
# -r 递归搜索; -n 显示行号; -H 输出文件名
# -A 打印出匹配行后面的几行; -B 打印前面几行; -C 打印前后几行
```