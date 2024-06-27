# 快速排序

main_test.go
    
```go
package main

import (
	"reflect"
	"testing"
)

func TestQuickSort(t *testing.T) {
	tests := []struct {
		input  []int
		output []int
	}{
		{
			input:  []int{4, 2, 1, 3},
			output: []int{1, 2, 3, 4},
		},
		{
			input:  []int{-1, 5, 3, 4, 0},
			output: []int{-1, 0, 3, 4, 5},
		},
		{
			input:  []int{},
			output: []int{},
		},
		{
			input:  []int{1},
			output: []int{1},
		},
		{
			input:  []int{10, 7, 8, 9, 1, 5},
			output: []int{1, 5, 7, 8, 9, 10},
		},
		{
			input:  []int{3, 3, 3},
			output: []int{3, 3, 3},
		},
	}

	for _, tt := range tests {
		nums := make([]int, len(tt.input))
		copy(nums, tt.input)
		QuickSort(nums)
		if !reflect.DeepEqual(nums, tt.output) {
			t.Errorf("QuickSort(%v) = %v, expected %v", tt.input, nums, tt.output)
		}
	}
}

```
    
main.go
    
```go
package main

func QuickSort(nums []int) {
	var (
		l = len(nums)
		i = 1
	)
	if l < 2 {
		return
	}
	for j := i; j < l; j++ {
		if nums[j] <= nums[0] {
			nums[i], nums[j] = nums[j], nums[i]
			i++
		}
	}
	nums[0], nums[i-1] = nums[i-1], nums[0]
	QuickSort(nums[:i-1])
	QuickSort(nums[i:])
}
```