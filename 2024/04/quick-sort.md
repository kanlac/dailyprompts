```go
package main

import "fmt"

func QuickSort(nums []int) {
	l := len(nums)
	if l <= 1 {
		return
	}

	i := 1
	for j := 1; j < l; j++ {
		if nums[j] < nums[0] {
			nums[i], nums[j] = nums[j], nums[i]
			i++
		}
	}
	nums[0], nums[i-1] = nums[i-1], nums[0]
	QuickSort(nums[:i-1])
	QuickSort(nums[i:])
}

func main() {
	s := []int{3, 6, 1, 0, 12, 4}
	QuickSort(s)
	fmt.Printf("%+v\n", s)
}
```
