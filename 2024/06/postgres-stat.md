# Postgres 数据统计/压力排查

## PostgreSQL 如何做压力排查

1. 针对连接数，定位到哪些客户端没有复用连接
2. 针对查询，定位到慢查询及其客户端

## 如何查询连接数

由于 Postgres 采用多进程架构，每创建一个连接都会产生一个进程，所以直接通过 ps -ef 也能大概判断连接数和客户端 IP，若需要更精确的数据则可以查询 pg_stat_activity 连接数据：

```sql
SELECT pid, usename, application_name, client_addr, client_port, state, query
FROM pg_stat_activity;

-- 统计各客户端的连接数
SELECT client_addr, count(*) AS connections
FROM pg_stat_activity
GROUP BY client_addr
ORDER BY connections DESC;

-- 查看最大连接数
SHOW max_connections;
```

如果连接数量达到上限，会出现 “too many clients” 错误。

## 如何定位慢查询

初次启用 pg_stat_statements 扩展，在 PostgreSQL 配置文件 postgresql.conf 中添加或取消注释以下行：

```
shared_preload_libraries = 'pg_stat_statements'
```

重启 PG 并登录，创建扩展：

```sql
CREATE EXTENSION pg_stat_statements;
```

查询：

```sql
-- 查看平均用时最长的 SQL
SELECT
    query,
    calls,
    total_time,
    mean_time,
    stddev_time,
    rows
FROM
    pg_stat_statements
ORDER BY
    mean_time DESC
LIMIT 10;

-- 修改 ORDER BY 字段为 calls，查看执行频率最高的 SQL
```

`pg_stat_statements` 的数据会在服务重启时重置，也可以通过 `pg_stat_statements_reset()` 函数手动重置。

refer http://www.postgres.cn/docs/9.5/pgstatstatements.html

## 如何定位慢查询机器对应的容器

组合查询：

1. 通过 pg_stat_statements 定位慢 SQL，提取 query 片段
2. 通过 query 片段和 pg_stat_activity 查到客户端 IP
    
    ```sql
    SELECT pid, usename, application_name, client_addr, client_port, state, query
    FROM pg_stat_activity WHERE client_port != -1 AND query LIKE '%SQL_SEGMENT%';
    ```
    
3. 根据 IP 查到容器名
    
    ```bash
    docker inspect --format='{{.Name}} - {{range .NetworkSettings.Networks}}{{.IPAddress}} {{end}}' $(docker ps -q)
    ```