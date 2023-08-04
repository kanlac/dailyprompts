# 数据库锁

### 何时使用乐观锁？如何使用？

乐观锁（Optimistic Locking）是一种在读取数据时不加锁，而在更新数据时才进行并发控制的策略。乐观锁假设数据在大部分时间内不会发生冲突，只有在实际更新时才会检查是否有冲突，因此被称为"乐观"。乐观锁适用于读多写少的并发场景。

以下是一个简单的乐观锁的使用示例：

1. **版本号机制**：这是乐观锁最常用的一种实现方式。在数据表中添加一个版本号字段，每次读取数据时，将版本号一同读出。在更新时，将此版本号一同提交，同时版本号加一。数据库根据提交的版本号和数据库中的版本号进行比较，如果数据库中的版本号大于提交的版本号，则拒绝更新；如果相等，则接受更新。

例如，在MySQL中，你可以这样使用版本号机制：

```sql
CREATE TABLE `user` (
  `id` int(11) NOT NULL,
  `username` varchar(255) DEFAULT NULL,
  `version` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

-- 当读取数据时，同时读取版本号
SELECT id, username, version FROM user WHERE id = 1;

-- 更新数据时，同时检查版本号，并且版本号加一
UPDATE user SET username = 'new_username', version = version + 1 WHERE id = 1 AND version = old_version;

```

在这个例子中，如果有两个并发的事务同时尝试更新同一条数据，只有一个可以成功，另一个会因为版本号不匹配而失败。

需要注意的是，乐观锁不能解决所有的并发问题，特别是在高并发的写入场景下，乐观锁可能会导致大量的更新冲突。在这种情况下，可能需要使用其他的并发控制策略，例如悲观锁或者分布式锁。

### 悲观锁有哪两种类型？

悲观锁（Pessimistic Locking）假定数据在任何时候都有可能被其他事务修改，因此在每次读取数据时都会加锁，确保在锁定期间，没有其他事务可以修改数据。这种策略适用于写操作较多的场景。

MySQL中的InnoDB存储引擎支持两种类型的悲观锁：共享锁（Shared Lock）和排他锁（Exclusive Lock）。

1. 共享锁（S锁）：多个事务可以同时对同一行数据加共享锁，一个事务在加共享锁之后，其他事务可以读该行数据但不能写（不能加X锁）。
2. 排他锁（X锁）：一次只允许一个事务对一行数据加排他锁，一个事务在加排他锁之后，其他事务不能对该行数据加任何锁。

以下是一个使用悲观锁的例子：

```sql
-- 启动一个新的事务
START TRANSACTION;

-- 使用共享锁读取数据
SELECT * FROM user WHERE id = 1 LOCK IN SHARE MODE;

-- 使用排他锁读取数据
SELECT * FROM user WHERE id = 1 FOR UPDATE;

-- 更新数据
UPDATE user SET username = 'new_username' WHERE id = 1;

-- 提交事务
COMMIT;

```

在这个例子中，`SELECT ... LOCK IN SHARE MODE`语句会给选中的行加共享锁，`SELECT ... FOR UPDATE`语句会给选中的行加排他锁。在事务提交或者回滚之前，其他事务不能修改被锁定的数据。

注意：在使用悲观锁时，务必要注意避免死锁。死锁是指两个或者多个事务在执行过程中，因争夺资源而造成的一种互相等待的现象，若无外力干涉，它们都将无法进行下去。一种常见的避免死锁的策略是总是以相同的顺序访问资源。