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

## 知识点

- 7 Queries
    - 7.8. WITH Queries (Common Table Expressions) refer https://dev.mysql.com/doc/refman/8.0/en/with.html
- 9 Functions and Operators
    - 9.9 时间处理：`EXTRACT` 提取时间中的某一部分，比如秒或月，PostgreSQL 还支持提取 EPOCH 数值。refer https://www.postgresql.org/docs/8.1/functions-datetime.html
    - 9.18 条件表达式：包括 `CASE WHEN`，`NULLIF` 等等，refer https://www.postgresql.org/docs/current/functions-conditional.html
    - 9.21 Aggregate function 聚合函数 - 使用多行，返回单行
    - 9.22 Window function 窗口函数 - 使用多行，返回多行。有 `OVER` 子句，OVER 子句定义了一个窗口或行集合，在这个窗口上，窗口函数进行计算；维基百科有比较不错的[介绍](https://en.wikipedia.org/wiki/Window_function_(SQL))
        - `LAG`：获取之前的记录；`LEAD`：获取之后的记录

## SQL 解读

- `SUM(is_new_group) OVER (ORDER BY time) AS group_id`
    
    `OVER (ORDER BY time)` 定义了窗口函数的操作范围。`OVER`子句定义了一个窗口或行集合，在这个窗口上，窗口函数进行计算。这里，`ORDER BY time`表示窗口中的行将根据`time`列的值进行排序。因此，对于结果集中的每一行，窗口函数都会计算该行及在该行之前（根据`time`排序）的所有行的`is_new_group`列的总和。
    
    因此，这一句SQL为结果集中的每一行计算了一个`group_id`，它代表了到当前行为止（根据`time`列排序），`is_new_group`列的累计总和。由于`is_new_group`列只包含0和1，所以每当`is_new_group`为1时，`group_id`就会增加，这实际上标识了一个新的组。这样，具有相同`group_id`值的行就会被认为属于同一个组。