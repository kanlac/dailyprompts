# 切片与数组

## Q1

```go
func editSlice(a []int) {
	a = append(a, 3)
	a[0] = 5
	fmt.Printf("inside: %+v\n", a)
}

func main() {
	s := make([]int, 0, 5)
	s = append(s, []int{1, 2}...)

	editSlice(s)
	fmt.Printf("outside: %+v\n", s)
}
```

问：输出是什么？

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

## 将指针 append 到切片，这个指针会发生拷贝吗？

会发生浅拷贝，也就是说会产生一个新的指针。

## 以下操作是否会报错

```go
s := []int{1}
s = s[1:]
```

answer:

不会报错。尽管用 1 作为下标是越界的，但做切片操作却不会。这是因为在 Go 语言中，切片操作使用的是半开区间 **`[start, end)`**，这意味着它包含起始索引，但不包含结束索引。

简单来说，当 **`start`** 和 **`end`** 相等时，结果是一个空切片，而不是错误。这也是为什么你的代码不会报“下标越界”的错误。这种设计允许 Go 程序员方便地进行切片操作，而不用担心边界条件，从而使得代码更加简洁和安全。