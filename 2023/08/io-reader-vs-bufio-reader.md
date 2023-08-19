# io.Reader v.s. bufio.Reader

## 什么区别？分别适用于什么场景？

`io.Reader` 和 `bufio.Reader` 都是 Go 语言中常用的读取接口和结构体，但它们在使用和场景上有所不同。

1. **io.Reader**:
    
    `io.Reader` 是一个接口，定义在 `io` 包中。它是 Go 中所有读操作的基础接口。它的定义很简单：
    
    ```go
    type Reader interface {
        Read(p []byte) (n int, err error)
    }
    
    ```
    
    **适用场景**:
    
    - 当你只是想从某个源中读取数据而不需要任何额外的缓冲或处理时。
    - 当你与其他期望 `io.Reader` 的 API（如很多标准库函数和方法）交互时。
    - 在需要实现自己的读取逻辑时，如自定义协议或格式的解码。
2. **bufio.Reader**:
    
    `bufio.Reader` 是一个结构体，定义在 `bufio` 包中。它提供了缓冲读取功能，这意味着它会从其底层 `io.Reader` 中提前读取并存储更多的数据，以减少实际的系统调用。此外，它还提供了很多实用的读取方法，例如 `ReadString`, `ReadByte`, `ReadRune` 等。
    
    **适用场景**:
    
    - 当你从一个涉及频繁系统调用的源（如文件或网络连接）读取数据时，使用 `bufio.Reader` 可以显著提高性能。
    - 当你需要按行读取文本文件时。`ReadString` 或 `ReadBytes` 使这变得很简单。
    - 当你需要更高级的读取方法，如读取一个单独的字节或 Unicode 码点。
    - 当你与文本数据交互并且需要更多的错误处理和控制时，如 CSV、JSON 或 XML 的解析。

**总结**:

`io.Reader` 提供了一个通用的读取接口，适用于所有读取需求。而 `bufio.Reader` 是一个针对特定场景的缓冲读取器，它可以提高读取效率，并提供更多的方法来简化文本数据的读取和处理。在设计和实现读取代码时，你应该根据具体的需求和场景来选择适当的读取器。

##  bufio.Reader 默认的缓存大小是多少？

4KB.