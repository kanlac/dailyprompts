# Postgres 自动清理

## 现象

每个月 1 号执行一个定时任务，删除某表中 30 天之前的数据，但发现删除后不能马上释放磁盘空间，因为该表的数据量非常大，因此开发人员在删除脚本中又加入了一个 VACUUM FULL 操作。

## 原理

1. VACUUM，只标记不清理，不获取锁，回收的空间不会释放，而是留给该表重用
2. VACUUM FULL，会获取一个排他锁（不推荐作为 routine 使用），重写表的全部内容到新的磁盘文件（需要额外的磁盘空间），因而能够释放尽可能多的空间到操作系统，操作慢
3. VACUUM TRUNCATE，会释放表最后的空白页的磁盘空间到操作系统，涉及排他锁
    - 官方推荐用自动清理 autovacuum，本质就是执行 VACUUM，但永远不会 VACUUM FULL），默认就是打开的

> Moderately-frequent standard `VACUUM` runs are a better approach than infrequent `VACUUM FULL` runs for maintaining heavily-updated tables.

- autovacuum 能够（根据更新操作）动态的调整清理周期，因此比自己设定一个固定时间清理要好，能够更好处理操作量激增的情况（否则还是得用 VACUUM FULL 应急）
- 不建议完全关闭 autovacuum，至少也要配置让它处理操作激增的情况
- 如果自己实现 vacuum，建议每天做一次数据库级别的 VACUUM（低使用量阶段）
- `vacuum_cost_xxx`：设定 I/O 上限，如果达到上限，隔一段时间再执行清理，不干扰其他数据库操作。这个设置也可以对手动执行的 VACUUM 命令开启
- 清理支持并发，可以配置多进程

总结：**建议的做法是提高删除的频率，而不是每月只删一次，然后执行 vacuum full。在正常的删除频率下，PG 自动清理是够用的**

```sql
-- 查看配置文件位置
SHOW config_file;
-- 查看配置的值
SHOW autovacuum;
-- 查看自动清理相关配置，以及配置来源（是默认还是配置文件）
SELECT name, setting, unit, short_desc, source FROM pg_settings where name like '%vacuum%';
-- 查看各个表自动清理执行的情况
SELECT relname, last_autovacuum, last_autoanalyze, autovacuum_count, autoanalyze_count FROM pg_stat_user_tables;
```
