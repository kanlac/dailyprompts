# Debug: 内存泄露排查

## 操作系统层面

top 或者 ps 的 `—-sort` 选项。

## Docker 层面

要分析哪些容器占用了较多的内存，您可以使用 `docker stats` 命令。以下是步骤：

1. 打开终端。
2. 输入 `docker stats` 命令。这会显示所有运行中的容器的实时统计信息，包括 CPU 使用率，内存使用量，网络 I/O，磁盘 I/O 等。
3. 查看 "MEM USAGE / LIMIT" 列。这列展示了每个容器的内存使用量和限制。如果内存使用量接近或达到了限制，那么这个容器可能就是内存使用高的容器。

注意：`docker stats` 默认会显示所有运行中的容器。如果只想看特定的容器，可以在命令后面加上容器的名字或 ID。

## K8s 层面

待补充。

## (Go) 程序层面

- 如何证明程序的堆内存占用情况？
    
    在 Go 语言中，要证明程序的堆内存占用情况，你可以使用几种方法：
    
    1. **`runtime` 包**：
    Go 的 `runtime` 包提供了多种函数来检查和报告内存使用情况。例如，`runtime.ReadMemStats` 函数可以用来获取各种内存统计数据，包括堆内存的使用情况。
        
        示例代码：
        
        ```go
        import (
            "runtime"
            "fmt"
        )
        
        func printMemStats() {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
            fmt.Printf("\\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
            fmt.Printf("\\tSys = %v MiB", bToMb(m.Sys))
            fmt.Printf("\\tNumGC = %v\\n", m.NumGC)
        }
        
        func bToMb(b uint64) uint64 {
            return b / 1024 / 1024
        }
        
        ```
        
    2. **分析工具（如 pprof）**：
    Go 提供了内建的分析工具，pprof，它可以用来分析程序的内存占用情况。pprof 可以生成堆内存的剖析文件，然后你可以使用 `go tool pprof` 命令来分析这个文件。
        
        使用 pprof 的步骤：
        
        - 在代码中导入 `net/http/pprof`。
        - 启动一个 HTTP 服务器来提供剖析接口。
        - 使用 `go tool pprof` 工具来获取和分析数据。
        
        示例代码：
        
        ```go
        import (
            _ "net/http/pprof"
            "net/http"
        )
        
        func main() {
            go func() {
                http.ListenAndServe("localhost:6060", nil)
            }()
            // 其他代码
        }
        
        ```
        
        然后，在终端中运行 `go tool pprof <http://localhost:6060/debug/pprof/heap`。>
        
    3. **使用第三方库**：
    也有一些第三方库，如 `github.com/pkg/profile`，可以帮助更容易地进行内存和CPU分析。
    
    通过使用这些方法，你可以有效地监控和证明 Go 程序的堆内存使用情况。
    
- 解释第一步获取的各项数据
    
    在 Go 语言中使用 `runtime.ReadMemStats` 获取的内存统计数据中，与堆内存相关的主要字段如下：
    
    1. **`Alloc`**:
        - 描述：这是当前堆上分配的字节数。
        - 重要性：反映了程序目前直接使用的堆内存量。
    2. **`TotalAlloc`**:
        - 描述：表示程序启动以来分配的总字节数（不考虑释放）。
        - 重要性：这个值随着时间的推移增加，即使有垃圾回收释放内存，它也不会减少。这有助于了解程序的总体内存分配趋势。
    3. **`Sys`**:
        - 描述：代表从操作系统中获得的内存总量，包括用于堆、栈和其他内部结构的内存。
        - 重要性：这反映了 Go 运行时从操作系统中获取的总内存量，但不仅限于堆内存。
    4. **`HeapAlloc`**:
        - 描述：与 `Alloc` 类似，表示当前在堆上分配的字节数。
        - 重要性：这是反映当前堆使用情况的关键指标。
    5. **`HeapSys`**:
        - 描述：表示操作系统分配给堆的内存总量，无论它是否正在使用。
        - 重要性：这包括了用于堆对象的内存以及未使用但未返还给操作系统的内存。
    6. **`HeapIdle`**:
        - 描述：堆上未使用且可以被操作系统回收的内存量。
        - 重要性：这部分内存是闲置的，表明了有多少堆内存目前没有用于存储任何对象。
    7. **`HeapInuse`**:
        - 描述：堆上当前正在使用或者保留未使用但不能被操作系统回收的内存量。
        - 重要性：这反映了活跃堆内存的使用情况。
    8. **`HeapReleased`**:
        - 描述：表示已经从堆中返回给操作系统的内存量。
        - 重要性：显示了 Go 运行时成功归还给操作系统多少内存，这有助于了解垃圾回收的效果。
    9. **`HeapObjects`**:
        - 描述：堆上当前分配的对象数。
        - 重要性：这可以帮助识别对象泄漏（如果这个数字不断增长并且没有相应的减少）。
    
    通过分析这些与堆内存相关的字段，可以更深入地了解 Go 程序的内存使用情况和模式。
    

### pprof

对于 Go 服务的内存分析，可以使用 Go 的内置 pprof 包进行动态内存分析。

1. 首先，在你的 Go 程序中导入 `net/http/pprof` 包。

```go
import _ "net/http/pprof"
```

1. 启动你的 Go 程序，确保它在运行中。
2. 使用 `go tool pprof <http://localhost:8080/debug/pprof/heap`> 获取 heap profile。这里的 `localhost:8080` 是你的服务的地址，可能需要根据实际情况调整。
3. 使用 `top` 命令查看内存使用最多的函数。或者使用 `web` 命令生成 SVG 格式的调用图。

注意：pprof 的数据可能会因为 Go 的垃圾回收而变动，所以你看到的结果可能会有波动。

### 基准测试

基准测试是测量代码性能的一种方法，它也可以用于测量代码的内存分配。在基准测试中，你可以使用`testing.B`对象的方法来开启和获取内存分配的信息：

```
func BenchmarkFoo(b *testing.B) {
    b.ReportAllocs() // 开启内存分配报告
    for i := 0; i < b.N; i++ {
        foo() // 测试foo函数的性能
    }
}

```

然后，你可以使用`go test -bench`命令来运行基准测试，并查看内存分配的结果。