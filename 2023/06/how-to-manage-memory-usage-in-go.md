# 如何做好 Go 的内存管理

### 安全篇

1. 在使用指针前总是判断是否为空
2. 使用切片的下标之前，必须进行长度校验
3. 注意整数安全｜运算有符号整数时，确保不会溢出；运算无符号整数时，确保不会反转

### 性能篇

基于 [性能问题的本质] 的第一个方面——不要做多余的事情——的注意事项：

1. 关于切片（一）：尽量不要将切片作为函数返回值，减少内存逃逸（[Go 的内存逃逸]），避免堆上垃圾过多
    
    ——但这不代表鼓励将切片作为参数传递！实际上，最好的做法仍然是将切片定义在下层（函数内部），并只返回只读 channel。以此实现 [隐式的并发安全设计]）
    
2. 关于切片（二）：若实在要将切片作为返回值，最好是拷贝一下，不要直接返回原来的切片，因为若返回的是一个很大的数组中的很小一份切片，空间浪费就会很严重
    - 具体实现
        1. 重置丢失的元素中的指针为 nil（参考[这里](https://gfw.go101.org/article/memory-leaking.html)）
        2. 使用 copy
            
            ```go
            func CopyDigits(filename string) []byte {
                b, _ := ioutil.ReadFile(filename)
                b = digitRegexp.Find(b)
                c := make([]byte, len(b))
                copy(c, b)
                return c
            }
            ```
            
        3. 使用 append
            
            ```go
            func CopyDigits(filename string) []byte {
                b, _ := ioutil.ReadFile(filename)
                b = digitRegexp.Find(b)
                return append([]byte{}, b...)
            }
            ```
            
3. 确保每个协程都能退出。每启动一个协程都相当于一次入栈操作，如果资源不进行回收，会导致**内存泄露**
4. 注意 [字符串处理]