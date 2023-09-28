# 编写一个时间聚合的 SQL

```sql
WITH SubQuery AS (
  SELECT
    probe,
    content,
    type,
    time,
    LAG(time) OVER (ORDER BY time) AS prev_time
  FROM
    events
),
GroupedData AS (
  SELECT
    probe,
    content,
    type,
    time,
    CASE WHEN EXTRACT(EPOCH FROM (time - prev_time)) > 100 THEN 1 ELSE 0 END AS is_new_group
  FROM
    SubQuery
),
AggregatedData AS (
  SELECT
    probe,
    content,
    type,
    time,
    SUM(is_new_group) OVER (ORDER BY time) AS group_id
  FROM
    GroupedData
)
SELECT
  probe,
  content,
  type,
  MIN(time) AS start_time,
  MAX(time) AS end_time,
  COUNT(*) AS aggregation_count
FROM
  AggregatedData
GROUP BY
  probe,
  content,
  type,
  group_id
ORDER BY
  start_time DESC;
```

- 7 Queries
    - 7.8. WITH Queries (Common Table Expressions) refer https://dev.mysql.com/doc/refman/8.0/en/with.html
- 9 Functions and Operators
    - 9.9 时间处理：EXTRACT 提取时间中的某一部分，比如秒或月，PostgreSQL 还支持提取 EPOCH 数值。refer https://www.postgresql.org/docs/8.1/functions-datetime.html
    - 9.18 条件表达式：包括 `CASE WHEN`，`NULLIF` 等等，refer https://www.postgresql.org/docs/current/functions-conditional.html
    - 9.21 Aggregate function 聚合函数 - 使用多行，返回单行
    - 9.22 Window function 窗口函数 - 使用多行，返回多行。有 `OVER` 子句；维基百科有比较不错的[介绍](https://en.wikipedia.org/wiki/Window_function_(SQL))
        - `LAG`：获取之前的记录；`LEAD`：获取之后的记录