# 随机数

## math/rand 与 crypto/rand 有何区别

前者生成可预测的随机数，速度更快，但如果攻击者知道了种子，就可以重现整个随机数序列。

后者是[加密安全的伪随机数生成器](https://en.wikipedia.org/wiki/Cryptographically_secure_pseudorandom_number_generator)，仍然用伪随机数，但是会用操作系统的熵池做种子，所以安全系数更高。比如 Linux 内核会维护一个熵池（entropy pool），用于存储随机数据的熵源。熵池收集来自多种系统事件的数据，例如键盘输入、鼠标移动、硬盘活动等，这些事件可以提供不可预测性。

真随机数的要求更严格，需要硬件采集严格意义上随机的信息，比如噪声、热波动等。

## Golang 如何高并发写入一个随机数数据

1. 确定长度，io.ReadFull 一次性写满整个数组
2. bufio Reader 避免在生成超长随机数时内存过度增长
3. 由于 bufio Reader 对象创建本身有开销，考虑用 sync.Pool
