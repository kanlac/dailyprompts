# kube-scheduler

Q:

- 说明一个 Pod 是如何完成调度的？
- 调度器是怎样做性能优化的？

两个调度算法：

1. 预选 Predicate，硬约束，计算符合 Pod 调度条件的节点
2. 优选 Priority，软约束，对可调度的节点打分

两个数据结构：

1. 调度器缓存 scheduler cache: 包括调度算法需要用到的节点和 Pod 的信息
2. 调度队列：存储待调度 Pod 的优先级队列，它又分为：
    1. activeQ，等待调度的 Pod
    2. backoffQ，临时性原因无法调度的、等待资源释放后可能可以调度的 Pod，会在退避时间后重试
    3. unschedulableQ，持久性原因无法调度的 Pod，需要等待集群状态变化后（比如新增节点）才会重新尝试调度

两个控制循环：

1. 通知路径 Informer Path: 监听 etcd 中 Pod 等 API 对象的变化，将待调度的 Pod 添加到调度队列
2. 调度路径 Scheduling Path: 从调度队列中出队一个 Pod，并从调度器缓存获取节点信息，用 Predicate 算法经过一遍过滤得到可以运行该 Pod 的节点列表，再用 Priority 算法对节点打分，并将 Pod 乐观绑定到得分最高的节点，并将 Pod 加入缓存

调度器性能优化三大核心：

1. 缓存化：调度器缓存用于将集群信息缓存化
2. 乐观绑定：绑定就是设置节点的 `nodeName`，绑定其实分成两步：1）Assume，调度路径的最后只假设调度完成，更新调度器缓存中的 Pod 和 Node 信息，不会阻塞，同时调度器会创建一个 goroutine 来异步地向 API Server 发起更新 Pod 的请求，更新绑定到 etcd；2）Admit，kubelet 二次确认该 Pod 能否在该节点运行。如果 Admit 操作中绑定失败了也没关系，等调度器缓存（在通知路径中）同步过来即可
3. 无锁化：调度器会启动多个 gorutine 以节点为粒度并发执行 Predicate 算法，用 MapReduce 并行计算然后汇总；Priority 算法也是一样。调度器会避免设置任何全局的竞争资源，免除了使用锁同步产生的巨大性能损耗。整个调度器只有在对调度队列和调度器缓存进行操作时才需要加锁，而这两个操作都不在调度路径上。补充：注意无锁化并不是整个调度过程完全没有用锁，实际上队列出队是要用锁的，但出队后马上就释放了，这里强调的是没有跨越整个调度过程的锁

![](kube-scheduler.webp)

源码笔记：

- Predicate 是怎么一个过程？是怎么搜寻节点的？——按索引一个个找，且会并发地找，找到一定个数就提前退出。每次调度会沿用之前的索引位置
- **Pod 的调度是不是并发的？——实际上，单个 Pod 的调度（SchedulingOne）是串行执行的，只是中间的调度算法是并发执行的，而且在完成计算后只更新缓存，上报 API Server 是异步完成的。串行处理的好处是避免竞态条件，而即便是需要并发的路径上，它也避免设置任何全局的竞争资源，免除了锁同步带来的巨大性能损耗**
- 队列出队时为什么没有使用 channel？——因为调度 Pod 需要一些时间，若把 Pod 放在 channel 里可能会过了时效，所以 Scheduler.NextPod 提供的是队列出队的方法，直接从队列中取。队列出队，没有用 channel，怎么实现阻塞呢？——用 sync.Cond，它是一种锁