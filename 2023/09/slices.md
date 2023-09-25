# 切片与数组

## Q1

Go（伪代码）：

```go
s := make([]int, 0, 5)
s = append(1, 2)

// 调用函数，传入 s
func editSlice(a []int) {
  a = append(3)
  a[0] = 5
}
```

问：最后 s 打印出来是多少？底层的数据又是多少？

答：s 是 5, 2，底层数组是 5, 2, 3, 0, 0

## 切片的结构？

可以把切片看作一个结构体，也称作 slice descriptor，其由 3 部分组成：

1. 指向底层数组某一个元素的指针

2. 切片的长度

3. 底层数组的容量（从指针位置开始）

```go
type slice struct {
    len int
    cap int
    array *array
}
```

## 作为函数参数时，是值拷贝还是指针拷贝？

是值拷贝。拷贝的是 slice descriptor，不过里边有指针，所以会拷贝指针。