# Golang 如何做错误处理

# 常用的错误处理包

1. 原生 `errors` 包
2. github/pkg/errors（以下用 `pkgerrors` 表示）

# 错误包装与解包

方案对比

1. 原生方案：`fmt.Errorf`，配合 `%w` 占位符做错误包装，`errors.Unwrap` 解包
    - 这样就完成了一次错误包装：`fmt.Errorf("run error: %w", err)`
    - `%w` 需要 Go 1.13+ 支持，否则需要自行实现包装器与解包器
    - 不推荐，写的时候很可能忘记包装
2. github/pkg/errors
    1. `pkgerrors.WithMessage` `pkgerrors.WithMessage` 仅含附加信息
    2. `pkgerrors.Wrap` `pkgerrors.Wrap` 含附加信息与栈信息
    3. `pkgerrors.WithStack` 仅含栈信息

解包

- 都有类似的解包实现： `errors.Unwrap` `pkgerrors.Unwrap`
- 都有 `Is` `As` 实现，可以自动解包错误链，比自行解包代码更简洁

# 自定义错误

如果我们需要对某种错误做识别比较操作，或者某一种错误可能出现多次，可以自定义错误。

### 静态错误

- Go 官方文档又称「哨兵（sentinel）错误」
- 定义哨兵错误，使用 `fmt.Errorf` 还是 `errors.New`？——推荐后者，因为哨兵错误都是静态信息，不需要动态格式化字符串
- Go 标准库已经定义了一些内置错误，比如我们可以用 `io.EOF` 识别文件结尾，可以用 `sql.ErrNoRows` 识别数据库查询没有返回结果的错误

### 自定义结构体

- error 本身的定义非常简单！它是一个接口类型，任何实现了 `Error() string` 方法的都可以作为 error
- 需要自行实现 `Error` `Unwrap` 方法

# Demo

```go
package main

import (
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

var sentinelErr = errors.New("a sentinel error")

func getWrappedError() error {
	return fmt.Errorf("get wrapped error: %w", sentinelErr)
}

func getWrappedPkgError() error {
	return pkgerrors.Wrap(sentinelErr, "get wrapped pkg error")
}

func getPkgErrorWithMessage() error {
	return pkgerrors.WithMessage(sentinelErr, "get a pkg error")
}

func getPkgErrorWithStack() error {
	return pkgerrors.WithStack(sentinelErr)
}

func main() {
	fmt.Println("===1===")
	fmt.Printf("%+v\n", getWrappedError())
	fmt.Println("===2===")
	fmt.Println(errors.Is(getWrappedError(), sentinelErr))
	fmt.Println("===3===")
	fmt.Printf("%v\n", getWrappedPkgError())
	fmt.Println("===4===")
	fmt.Printf("%+v\n", getWrappedPkgError())
	fmt.Println("===5===")
	fmt.Println(pkgerrors.Is(getWrappedPkgError(), sentinelErr))
	fmt.Println("===6===")
	fmt.Printf("%+v\n", getPkgErrorWithMessage())
	fmt.Println("===7===")
	fmt.Printf("%+v\n", getPkgErrorWithStack())
}
```