# Debug: 内存泄露排查

## 操作系统层面

top 或者 ps 的 `—-sort` 选项。

## 容器层面

要分析哪些容器占用了较多的内存，您可以使用 `docker stats` 命令。它显示所有**运行中的容器**的实时统计信息，包括 CPU 使用率，内存使用量，网络 I/O，磁盘 I/O 等。

查看 "MEM USAGE / LIMIT" 列。这列展示了每个容器的内存使用量和限制。如果内存使用量接近或达到了限制，那么这个容器可能就是内存使用高的容器。

注意：`docker stats` 默认会显示所有运行中的容器。如果只想看特定的容器，可以在命令后面加上容器的名字或 ID。

## Go 程序层面

在 Go 语言中，要证明程序的堆内存占用情况，你可以使用几种方法：

1. 标准库
    
    Go 的 `runtime` 包提供了多种函数来检查和报告内存使用情况。例如，`runtime.ReadMemStats` 函数可以用来获取各种内存统计数据，包括堆内存的使用情况。
    
    `runtime.ReadMemStats` 可以获取一些基本的内存分配相关的数据，**比如内存总量，活跃的堆内存量（HeapInuse），闲置的堆内存量（HeapIdle）等等**
    
    - 示例代码
        
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
        
2. 分析工具
    - pprof
        
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
        
3. 基准测试
    
    使用`testing.B`对象的方法来开启和获取内存分配的信息：
    
    ```
    func BenchmarkFoo(b *testing.B) {
        b.ReportAllocs() // 开启内存分配报告
        for i := 0; i < b.N; i++ {
            foo() // 测试foo函数的性能
        }
    }
    ```
    
    然后，你可以使用`go test --bench`命令来运行基准测试，并查看内存分配的结果。
    
4. 第三方工具
    
    如 `github.com/pkg/profile`，可以帮助更容易地进行内存和CPU分析。