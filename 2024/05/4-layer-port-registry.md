# 四层端口注册服务

- 设计方案——包括如何实现动态增删，包括要给服务配置 cluster role rbac
- 原理——验证出一套 ingress-nginx 开启四层端口的方案
- 困难——一致性哈希，因为 SNAT 会导致丢失源 IP，使 service 配置的亲和性起不了实际作用，方案比较： 1）`externalTrafficPolicy: Local`（用于外部均衡器）；2）加 ingress 标注（仅适用于七层）；3）编写 lua 脚本基于节点 ip 做哈希
