# 发起 DNS 请求

# Q
- DNS 端口号是多少？
- 为什么使用 UDP 而不是 TCP？
- 简述 DNS 请求的完整过程
- 如何手动进行域名解析？

# A
DNS 端口号是多少？——53。

为什么使用 UDP 而不是 TCP？
1. UDP 更快，不需要 3 次握手
2. DNS 请求的数据量很小，适合放到 UDP 数据报里发送
3. 尽管 UDP 本身不可靠，但我们可以在应用层为其实现可靠性，比如超时处理

简述 DNS 请求的完整过程——分为无缓存与有缓存两种情况……

如何手动进行域名解析？——使用 nslookup 可以直接与 DNS 服务器交互，而且可以直接向权威服务器发起查询
