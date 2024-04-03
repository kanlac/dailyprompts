# 用户登录/身份验证/授权

## 用户登录怎么做？不同方案的区别是什么？

1. Basic Auth：一种简单的**无状态**的用户名+密码认证。需要在**每次** HTTP 请求的 **Authorization** 字段中提供（base64 编码的）用户名和密码
2. 基于会话（Session）的认证：不用在每个请求中的 Authorization 包含用户名和密码
3. OAuth：OAuth 是一个开放标准，允许用户授权第三方网站访问他们存储在另外的服务提供者上的信息，而不需要给出密码。例如，“使用 Google 账户登录”就是基于 OAuth 的

## Basic Auth 

一种简单的**无状态**的用户名+密码认证。需要在**每次** HTTP 请求的 **Authorization** 字段中提供（base64 编码的）用户名和密码。

如果在返回的 header 中附加上 `WWW-Authenticate` 字段，一般浏览器会自动识别并跳出输入框。

注意：

- 是 base64 编码，但没有加密，基本等于明文传输，所以一定要搭配 HTTPS/TLS 使用

Basic Auth 怎么做成有状态的？——Basic Auth 是无状态的，如果想要实现有状态的登录，比如说设定超时时间，可以只在第一次认证的时候用 basic auth，认证通过后创建一个会话。

Nginx 也支持配置 Basic Auth，使用`htpasswd`工具来创建一个密码文件后，可以在你想保护的 location 中使用它。
