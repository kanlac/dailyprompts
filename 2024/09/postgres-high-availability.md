# Postgres 高可用

## Postgres HA 需要考虑什么问题

1. 如何读取
2. 如何写入
    1. 几个节点可写？主备是一种常用的方案，只有主节点可写，备用节点分为 hot standby（可读的实时备份）和 warm standby（延时备份，直到成为主节点之前不可读）
    2. 同步写入保证高一致性，还是异步写入保证高可用？同步写入代表只有所有节点完成写入事务才结束，故障转移时不会丢失数据
3. 节点故障怎么处理？如何做故障转移
    - 故障转移后，备变主，新主需要在旧主起来后通知其成为备，Postgres 本身不提供这样的机制
    - 旧主起来后，可以使用 pg_rewind 工具修复时间线分裂，同步新主上的变更数据
    - 故障转移后，可以是等待旧主恢复，也可以是在第三个节点上创建新备

## Postgres HA 可以怎么做

https://www.postgresql.org/docs/current/high-availability.html

1. 存储级别
    1. 共享磁盘，NAS，如果一个节点挂了，另一个节点调度新的服务并挂载这个磁盘，优势是不需要处理数据同步，快速的故障转移
    2. File System (Block Device) Replication，在文件系统层面做备份，做一个文件系统镜像，参考 DRBD 实现
    
    官方文档不推荐在存储级别做同步，应用级别 Postgres 已经提供了可靠的 WAL 同步，原因参见 https://github.com/cloudnative-pg/cloudnative-pg/blob/main/docs/src/architecture.md#synchronizing-the-state。
    
2. 应用级别
    1. 物理复制 - WAL，主流的高可用方案都基于此实现
        - 写日志 → 提交事务 → 发送日志
        - 异步，事务提交后发送日志，因此有数据丢失风险
        - 在异步的基础上，又分为
            1. 大颗粒度的文件同步，温备，非 HA
            2. 小颗粒度的流复制，热备，HA
    2. 逻辑复制 - 复制变化（SQL）而非字节；因为逻辑更改与数据库实现无关，所以可以跨主要版本复制，比如从 PG 11 到 PG 13，也可以跨不同操作系统或硬件架构复制；因为需要解析 SQL，性能开销较大，不适用于故障切换；不支持多主复制或双向复制
3. 多主复制
    1. 异步多主复制，适用于间歇性的同步，例如笔记本电脑和服务器
    2. 同步多主复制，优势是每个节点都可以写，无需区分主备，劣势是写入性能差，不适合以读为主的场景，但可以搭配共享磁盘

## 主流 HA 应用

1. Citus，PG 扩展，将 PostgreSQL 转变为分布式数据库，通过 coordinator 节点完成查询分发，但本身不提供高可用和故障转移机制，需要搭配 Patroni 或其他工具
2. Patroni - PG 高可用模版，故障转移管理工具

## K8S 中的主流 HA 应用

基于 DCS (distributed control system) 或 Kubernetes Operator 的方案，operator 通过声明式方式实现故障转移等自动化运维管理。无奇数个节点要求，只要底层的一致性存储（控制节点数量）是奇数。 [分布式共识](https://www.notion.so/41365e1f13b64aaab33bfdf688c7cce5?pvs=21) 

1. CrunchyData https://github.com/CrunchyData/postgres-operator, 3.8k
2. Zalando https://github.com/zalando/postgres-operator, 4.1k，简化 k8s 中 patroni 的使用，包含 k8s 资源和角色的创建，使用 operator 做编排工作
    - 通过定义 `postgresql` 资源创建 pg 集群，实际跑的容器是 [spilo](https://github.com/zalando/spilo?tab=readme-ov-file)
    - spilo 使用 WAL-E 和 [WAL-G](https://github.com/wal-g/wal-g) 做 WAL 的持续归档
    - 仅提供一个 Service
    - 本身不提供监控，参考文档 sidecar 部分了解如何配置 prometheus
    - 有[人](https://news.ycombinator.com/item?id=37618997)反应是 Zalando 内部使用的开源项目，对更广泛的使用场景支持不太好，可能是个缺陷
    - 如何从外部连接 PG
        
        ```bash
        export PGMASTER=$(kubectl get pods -o jsonpath={.items..metadata.name} -l application=spilo,cluster-name=acid-minimal-cluster,spilo-role=master -n zalando-pg)
        export PGREPLICA=$(kubectl get pods -o jsonpath={.items..metadata.name} -l application=spilo,cluster-name=acid-minimal-cluster,spilo-role=replica -n zalando-pg)
        kubectl port-forward $PGMASTER 6432:5432 -n zalando-pg
        
        export PGHOST=127.0.0.1
        export PGPORT=6432
        export PGUSER=postgres
        export PGPASSWORD=$(kubectl get secret -n zalando-pg postgres.acid-minimal-cluster.credentials.postgresql.acid.zalan.do -o 'jsonpath={.data.password}' | base64 -d)
        export PGSSLMODE=require
        psql
        ```
        
3. https://github.com/cloudnative-pg/cloudnative-pg, 3.6k, 文档翔实，[PPT](https://www.slideshare.net/slideshow/cloud-native-postgresql/243775892)，基于 CRD
    - 综合反馈最好
    - 声明式，使用简单
    - 不依赖 Patroni 等故障转移工具
    - 通过定义 `Cluster` 资源创建 pg 集群，使用镜像 `ghcr.io/cloudnative-pg/postgresql:13.6`
    - 提供 kubectl 插件，与集群交互很方便
    - 遵循 Kubernetes 的标准，设置较为规范
    - 使用 [barman](https://pgbarman.org/) 做持续归档/备份
    - 提供 3 个 Service：指向主节点的 rw，指向热备只读的 ro，指向所有只读的 r。故障转移后 cnpg 会自动更新 service
    - 高可用：除非整个可用区（即 Kubernetes 集群的所有节点）都出现故障，否则 CloudNativePG 集群不会发生故障
    - shared-nothing 架构，除了网络什么都不共享，存储使用本地卷