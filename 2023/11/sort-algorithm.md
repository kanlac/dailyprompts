# 排序算法

时间复杂度为 O(n²) 的有冒泡排序，选择排序，插入排序，采用分治法的快速排序效率最高，时间复杂度仅为 O(nlogn)。

快速排序实现：

```go
package main

import "fmt"

func quickSort(s []int) {
	var (
		l = len(s)
		i = 1 // i 下标左侧的数都小于 pivot
	)
	if l < 2 {
		return
	}

	pivot := s[0]
	for j := 2; j < l; j++ {
		if s[j] < pivot {
			s[i], s[j] = s[j], s[i]
			i++
		}
	}
	s[0], s[i-1] = s[i-1], s[0]

	quickSort(s[:i-1])
	quickSort(s[i:])
}

func main() {
	s := []int{3, 6, 1, 0, 12, 4}
	quickSort(s)
	fmt.Printf("%+v\n", s)
}
```