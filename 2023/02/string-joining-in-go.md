### Go 中字符串的拼接方式有哪些？各自效率如何？

字符串拼接的方法：
1. `+` 和 `+=` 操作符
2. `bytes.Buffer` 的 `WriteString()` 方法
3. `fmt.Sprintf()` 方法
4. `strings` 包中，`Join()` 方法或者 `Builder` 的 `WriteString()` 方法
5. 使用 `[]byte`

性能分析：
- `fmt` 要用到反射，所以拼接效率最低，`+` 其次；其他几个要好很多，差别没那么大，但性能最好的是预设容量的 byte 拼接
- 综合易用性和性能，一般推荐使用 `strings.Builder` 来拼接字符串