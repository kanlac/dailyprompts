# 一致性哈希

- 一致性哈希解决什么问题？
- 如何解决分布不均衡的问题？

一致性哈希解决什么问题？在水平扩展领域中，当集群中节点数量发生变化，不仅仅是挂掉的节点相应的连接要重置，而是所有连接都要重置，导致缓存丢失风暴（以缓存服务为例）。一致性哈希要解决的问题就是当节点数量变化时（也就是哈希表发生 resize 时），只有少量的键需要重新分布。

画一个（没有取模运算的）哈希环，哈希结果（服务器）都落在这个环上。对 key 做哈希运算，结果落在环上的某一点，走顺时针，取遇到的第一个服务器作为哈希结果。这样，不管是增加还是删除节点，都只有哈希环上的一部分映射需要重新处理。这就是一致性哈希算法。

一致性哈希存在的问题：因为要考虑增加和移除节点，无法确保所有服务器在环上占有的分区大小相等，所以计算出某个服务器的概率可能相当大，另一个则相当小。

为了解决一致性哈希分布不均衡问题，可以使用虚拟节点/副本技术。每个真实节点由一定数量的虚拟节点替代，将这些虚拟节点放在哈希环上。虚拟节点数量越多，分布就越均衡。