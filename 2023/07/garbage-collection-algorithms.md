# 常见垃圾回收算法

### 引用计数 Reference Counting

每个对象有个计数器，维护自己被引用的次数，当计数归零时，认为是不可达，将其进行回收。

- 优点：能立即回收不再被引用的对象
- 缺点：频繁更新计数会带来很多开销；无法解决循环引用的问题

### 标记清除 Mark and Sweep

这种算法分为两个阶段。在标记阶段，它遍历所有活动对象并将它们标记为“活动”。在清除阶段，它遍历所有对象，并删除未被标记为“活动”的对象。

- 优点：简单直接
- 缺点：首先，标记和清除阶段都需要停止程序执行（stop-the-world），这会导致程序的执行性能下降。其次，清除阶段会导致内存碎片化，这可能会影响后续的内存分配效率

### 标记整理 Mark and Compact

与标记-清除（Mark-Sweep）算法类似，标记-整理算法的第一步也是标记阶段，这一阶段会遍历所有的活动对象，并将它们标记为“活动”。

但在第二步，即整理阶段，标记-整理算法与标记-清除算法有所不同。整理阶段会将所有的活动对象移动到内存的一端，然后直接回收剩余的内存，而不是单独清除每个死亡对象。

- 优点：解决了标记清除法内存碎片化的问题
- 缺点：首先，移动对象需要花费额外的时间，尤其当活动对象数量很大时。其次，如果活动对象包含对其他对象的引用，那么在移动对象时，还需要更新这些引用，这会使算法更加复杂

### 复制法 Copying

此算法将内存分为两个等大的区域，每次只使用其中一个。当这个区域满了，算法就会找出所有活动的对象，然后把它们复制到另一个区域，并把原区域清空。

- 优点：解决了标记清除法内存碎片化的问题
- 缺点：需要一半的内存始终处于空闲状态，这意味着内存利用率较低

### 分代法 Generation

此算法的基本思想是将所有的对象分为几代，并且基于经验统计，假定新生代的对象会更快地变为不可达。因此，它频繁地收集新生代，但只偶尔收集老年代。

- 优点：减少了GC的总体开销
- 缺点：对于生命周期分布不均的应用效果较差。分代算法的效率基于一个假设，即大部分对象的生存周期都较短。但对于某些特殊的应用，例如大量使用了缓存，或者对象的生命周期分布更均匀，这个假设就不再成立，分代收集的效果就可能不如预期

### 并发 Concurrent

此算法允许垃圾收集和应用程序并发运行，而不是停止整个程序进行垃圾收集。与增量GC类似，其目标是减少应用程序的停顿时间。

- 优点：减少停顿时间
- 缺点：实现更复杂，需要更多的协调开销

### ……