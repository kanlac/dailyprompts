# 只出现一次的数字 II

https://leetcode.cn/problems/single-number-ii/description/

状态机解法相当有挑战，需要一步一个脚印地理解。

第一个关键点在于理解「**遇 1 进位就是做异或**」，记住这个定律。也就是说 one 需要遇 1 进位时有 `one ^= n`。

第二个关键点在于通过状态机**看到什么时候 one 不需要遇 1 进位**。

根据 00 -> 01 -> 10 的状态图发现，two 的当前值为 1 时，one 不需要遇 1 进位，于是可以得到：

```
if two == 0:
    one ^= n
if two == 1:
    one = 0
```

第三个关键点在于**看到什么时候 two 不需要遇 1 进位**。one 计算完若为 1，不论它是否经过进位，two 一定为 0（因为最高就是 01 了没有 11）；其它情况下，two 正常遇 1 进位。

```
if one == 1:
    two = 0
if one == 0:
    two ^= n
```

写出来发现和 two 和 one 的计算方式是一样的。

第四个关键点在于把上面的抽象表达式转化为实际可用的位运算，得：

```go
ones = ones ^ num & ~twos
twos = twos ^ num & ~ones
```
