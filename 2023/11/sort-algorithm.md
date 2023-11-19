# 排序算法

时间复杂度为 O(n²) 的有冒泡排序，选择排序，插入排序，采用分治法的快速排序效率最高，时间复杂度仅为 O(nlogn)。

快速排序实现：

```go
package main

import "fmt"

func quickSort(s []int) {
	if len(s) <= 1 {
		return
	}
	var (
		pivot = s[0]
		x     = 1 // x 下标左边的数都是小于 pivot 的
		l     = len(s)
	)
	for i := 1; i < l; i++ {
		if s[i] < pivot {
			s[i], s[x] = s[x], s[i]
			x++
		}
	}
	s[0], s[x-1] = s[x-1], s[0]
	quickSort(s[:x-1])
	if x < l {
		quickSort(s[x:])
	}
}

func main() {
	s := []int{3, 6, 1, 0, 12, 4}
	quickSort(s)
	fmt.Printf("%+v\n", s)
}
```

插入排序实现：

```go
package main

import "fmt"

func InsertSort(nums []int) {
	l := len(nums)
	if l < 2 {
		return
	}

	for i := 1; i < l; i++ {
		j := i
		cur := nums[i]
		for j > 0 && nums[j-1] > cur {
			nums[j] = nums[j-1]
			j--
		}
		nums[j] = cur
	}
}

func main() {
	s := []int{3, 6, 1, 0, 12, 4}
	InsertSort(s)
	fmt.Printf("%+v\n", s)
}
```