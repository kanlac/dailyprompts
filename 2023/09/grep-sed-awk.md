# Q

- 如何执行目录下的文本搜索，同时打印附近的行与行号？
- 演示 sed 如何给某行加上注释，或者删除某行

# grep

```bash
# 筛选打印日志中的所有包含 "ERROR" 的行，倒序（inverting）
grep -v "ERROR" log.txt

# 如何执行目录下的文本搜索，同时打印附近的行与行号？
grep -rnH -C 2 "foo" .
# -r 递归搜索; -n 显示行号; -H 输出文件名
# -A 打印出匹配行后面的几行; -B 打印前面几行; -C 打印前后几行

# 指定层级进行文本搜索
find . -maxdepth 2 -type f -exec grep -H "your-pattern" {} \;
```

# sed

流编辑器，常用于替换文本。

```bash
# Usage: sed [OPTIONS] SCRIPT FILE...
# SCRIPT 是一个这样结构的字符串：
#   [addr]X[options]
#   [addr] 是条件过滤，或者叫地址选择器，X 是要执行的 sed 命令
```

## 删除

`/模式/d`

```bash
# 删除包含 foo 的行
# -i 直接修改原文件而非输出
sed -i '/foo/d' filename

# 删除 # 号开头的行
sed '/^#/d' filename

# 删除第五行
sed '5d' filename

# 删除 2-4 行
sed '2,4d' filename
```

## 替换

`s/原模式/新模式/标志`

- 标志
    
    精准控制，如果不提供，模式是只会替换每一行中第一个匹配项
    
    - `g`：全局替换，即替换所有匹配项
    - `i`：匹配时忽略大小写
    - 数字：仅替换第 n 个匹配项

```bash
# 将每行第一个 apple 替换为 orange
sed 's/apple/orange/' file.txt

# 将所有 apple 替换为 orange
sed 's/apple/orange/g' file.txt

# 将每行第二个 apple 替换为 orange
sed 's/apple/orange/2' file.txt

# 注释掉包含 'tcp_tw_recycle' 的行
# /tcp_tw_recycle/ 模式匹配，定位包含 'tcp_tw_recycle' 的行
# s/^/#/ 替换命令，将匹配行的开始处替换为 '#'
sed -i '/tcp_tw_recycle/ s/^/#/' /etc/sysctl.conf
# 这里地址选择器与命令之间有一个空格，它可有可无
```

# awk

awk 是 3 中工具里最强的，它实际上是一种用于操作数据、生成报告的脚本语言。在 awk 语言中，有变量，有逻辑运算，也有函数。

awk 对输入文本按行扫描，进行模式匹配（prog）
