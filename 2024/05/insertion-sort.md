# 插入排序

main_test.go

```go
package main

import (
	"reflect"
	"testing"
)

func TestInsertionSort(t *testing.T) {
	tests := []struct {
		name     string // 测试用例的名称
		input    []int  // 输入数组
		expected []int  // 预期的排序后数组
	}{
		{"sorted list", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reverse list", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"empty list", []int{}, []int{}},
		{"single element", []int{1}, []int{1}},
		{"random list", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InsertionSort(tt.input)
			if !reflect.DeepEqual(tt.input, tt.expected) {
				t.Errorf("InsertionSort(%v) got %v, want %v", tt.input, tt.input, tt.expected)
			}
		})
	}
}
```

main.go

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

