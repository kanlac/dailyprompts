# Go 接口实现

### 问题

以下代码输出什么？为什么？

一、
```go
// chapter5/sources/interface-internal-1.go

type MyError struct {
    error
}

var ErrBad = MyError{
    error: errors.New("bad error"),
}

func bad() bool {
    return false
}

func returnsError() error {
    var p *MyError = nil
    if bad() {
        p = &ErrBad
    }
    return p
}

func main() {
    e := returnsError()
    if e != nil {
        fmt.Printf("error: %+v\n", e)
        return
    }
    fmt.Println("ok")
}
```

二、
```go
func printEmptyInterfaceAndNonEmptyInterface() {
    var eif interface{} = T(5)
    var err error = T(5)
    println("eif:", eif)
    println("err:", err)
    println("eif = err:", eif == err)
}
```

### 基础知识

- Go 语言中每种类型都有唯一的 `_type` 信息，无论是内置类型还是自定义类型。

### 内部实现

Go 的接口变量有两种，一种是有方法的 `iface`，一种是没有方法的 `eface`。它们的内部结构是这样：

```go
// $GOROOT/src/runtime/runtime2.go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}

type eface struct {
    _type *_type
    data  unsafe.Pointer
}
```

第二个字段 `data` 的功用相同，都指向当前赋值给该接口类型变量的动态类型变量的值。

接下来看 `tab` 和 `eface` 的差别：

```go
// $GOROOT/src/runtime/runtime2.go

type itab struct {
    inter *interfacetype
    _type *_type
    hash  uint32
    _     [4]byte
    fun   [1]uintptr
}

// $GOROOT/src/runtime/type.go

type _type struct {
    size       uintptr
    ptrdata    uintptr
    hash       uint32
    tflag      tflag
    align      uint8
    fieldalign uint8
    kind       uint8
    alg        *typeAlg
    gcdata    *byte
    str       nameOff
    ptrToThis typeOff
}
```

后者存储的是接口的类型信息（如前面所提示的，Go 中每个类型都有自己的 `_type`），而前者不仅包含了后者，还额外有方法列表。具体来说， `fun` 是动态类型已实现的接口方法的调用地址，而 `inter` 存储着该接口类型自身的信息：

```go
// $GOROOT/src/runtime/type.go
type interfacetype struct {
    typ     _type
    pkgpath name
    mhdr    []imethod
}
```

其中，`mhdr` 是接口方法的集合。

### 举例与映证

如果声名一个接口类型并且不初始化，那么它就是一个 `eface`，并且 `_type` 和 `data` 都为空，这时候是可以正常判空的。

返回到开头的那个问题，`returnsError` 返回的时候会把变量 `p` 封装成 `error` 接口类型变量，而我们知道作为一个内置的接口，`error` 是有方法的，所以这时候它是一个 `iface`，且接口信息表（`itab`）并不为空，因此即便 `data` 为空，它也不是 `nil`(0x0,0x0)。

### 两种接口的等值比较

针对第二个问题做出解答。

前提知识：看似空接口类型与非空接口类型是无法相等的，因为 `_type` 和 `itab` 不会一样，但实际上 Go 在对非空接口做等值比较的时候用的是 `itab._type` ，所以是可以相等的。

程序输出：

```
eif: (0x10b3b00,0x10eb4d0)
err: (0x10ed380,0x10eb4d8)
eif = err: true
```

做等值判断的时候，即时接口类型不一样也没有关系，比较的是动态类型的值。