# 如何减少锁的使用？

减少数据库锁使用的常用方法，逐层递进：

1. 查什么——**合理使用索引**：利用索引能够减少**全表扫描**，相应的，也就减少了锁定的数据范围
2. 怎么查——**合理的并发控制**：控制系统的并发级别，防止过度并发导致大量的锁等待和竞争。
3. 查多少——**优化查询，减少事务大小**：尽可能使每个事务的操作尽量小且简短，这样就可以减小锁粒度，减少锁定的时间和范围
4. 什么锁——**使用乐观锁**：如果业务允许，可以使用乐观锁，它通过在数据更新时检查数据是否被修改来减少锁的使用
5. **合理设置事务隔离级别**：在保证数据一致性的前提下，选择适当的事务隔离级别可以减少不必要的锁
