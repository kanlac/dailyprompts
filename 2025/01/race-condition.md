# 竞态问题

## 解决竞态问题的常见方法

### 共享内存

- 互斥锁：种类繁多且复杂，包括自旋锁、悲观锁和乐观锁
- 原子操作：线程安全而且性能很高，因为它是从 CPU 硬件层面做的更新，使用 [CAS (Compare And Swap)](https://www.notion.so/CAS-Compare-And-Swap-d333ca0650034f9a9462ac055e3bd52c?pvs=21) 指令来更新原子变量。
- STM：相对于硬件/数据库层面的并发控制机制，software transactional memory 是指在软件层面实现关系型数据库中的事务。一个事物必须满足 ACID 性质。

### 消息通信

使用信号量，消息传递而不是共享内存。相比锁模型更安全高效。包括 CSP 模型和适用于分布式场景的 Actor 模型。

相关输出： [隐式的并发安全设计](https://www.notion.so/2411a18da6134d7d8c06651334ff6367?pvs=21) 

### 写时复制

另一种无锁机制，Copy on write (COW)，Docker 构建镜像就通过它来实现快速读写。

ref: 腾讯云张翔的[分享](https://mp.weixin.qq.com/s/EFi1GzHy5qAx9Ixnppoybw)

## 为什么 CoW 能实现瞬间复制，还能避免冲突

refer blog: https://kanlac.in/how-to-speed-up-ci-pipeline
