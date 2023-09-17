# Go 垃圾回收

## Go 的垃圾回收器使用了哪种算法？

Go 采用的垃圾回收算法是标记清除（refer [常见垃圾回收算法](https://www.notion.so/0537956cf06f4b8da12d143577cf3522?pvs=21)）。在标记阶段，GC会从根对象（栈上的对象、全局对象）开始，遍历所有可达对象，并进行标记。在清除阶段，清除所有未被标记的对象。

## Go语言的垃圾回收器在什么时候会开始工作？有什么触发条件？

Go语言的垃圾回收器（GC）采用了一种名为"并发标记扫描"（Concurrent Mark Sweep, CMS）的机制。在这种机制下，Go的GC主要会在以下条件下启动：

1. 堆内存分配量达到一定阈值：如前所述，垃圾收集器会在堆内存增长到设定的阈值时启动新的收集。这个阈值由 `GOGC` 环境变量控制。默认情况下，`GOGC` 的值是 `100`，表示当堆内存增长到上次GC后的两倍时，就会触发新的垃圾收集。
    - 修改 GOGC
        
        `GOGC=800 ./your-program`
        
        这会让程序使用更多内存，运行更快
        
2. 显式调用：如果程序明确调用了 `runtime.GC()` 函数，那么垃圾回收器将立即启动。
3. 系统内存压力：当系统内存压力过大，可能会触发Go运行时启动垃圾收集。

需要注意的是，尽管存在上述触发条件，但具体的GC启动时间并不是严格固定的。Go的运行时系统会根据实际情况（包括CPU使用情况，内存使用情况等）来选择合适的时机启动垃圾收集，以尽量减少对程序运行的影响。

## 垃圾收集的四个阶段

来自最新版本 1.21.1 源码的注释 [https://github.com/golang/go/blob/2c1e5b05fe39fc5e6c730dd60e82946b8e67c6ba/src/runtime/mgc.go#L24）：](https://github.com/golang/go/blob/2c1e5b05fe39fc5e6c730dd60e82946b8e67c6ba/src/runtime/mgc.go#L24%EF%BC%89%EF%BC%9A)

1. 第一阶段 sweep termination，清理终止，会触发 STW，使得所有 P 达到安全点，并且会清除（sweep）任何尚未清理的内存段，除非这个 GC 周期是在预期时间之前强制执行的，否则这一阶段结束的时候就不会有未被清理的内存段了
2. 第二阶段 mark phase，标记阶段
    1. 开始入队标记任务（mark worker），同时**会有一个短暂的 Stop the World 确保所有 P 的写屏障已开启**
    2. 写屏障都开启后，**Start the World，开始执行标记任务**（mark worker）。也就是说，并发的标记任务和程序的 goroutine 在同时执行
    3. 执行根节点标记，包括扫描所有的栈，所有的全局变量，扫描一个 goroutine 的时候，会短暂地停止该 goroutine 以获取其堆栈上的所有指针
    4. ……
    5. gc 是分散式、并发式进行的，当没有 root marking job 或者灰色对象时，GC 会进入下一阶段
3. 第三阶段 mark termination，标记终止，会再次触发短暂的 STW
4. 第四阶段 sweep phase，进入清理阶段，关闭写屏障，start the world 恢复 goroutine，从此时开始，新分配的对象都是白色，在后台并发执行清理操作
5. 进行了足够的内存分配（达到 GC rate）后，重复以上过程

细节解释：

- 写屏障（write barrier），是用于追踪内存写操作的机制，在 start the world 之后，新分配的对象会被标记为黑色，表示不会在此次 GC 周期中被回收
