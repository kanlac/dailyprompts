# 三数之和

https://leetcode.cn/problems/3sum/description/

main_test.go

```go
package main

import (
	"reflect"
	"testing"
)

func TestThreeSum(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected [][]int
	}{
		{
			name:     "basic example",
			nums:     []int{-1, 0, 1, 2, -1, -4},
			expected: [][]int{{-1, -1, 2}, {-1, 0, 1}},
		},
		{
			name:     "empty array",
			nums:     []int{},
			expected: nil,
		},
		{
			name:     "all zeroes",
			nums:     []int{0, 0, 0},
			expected: [][]int{{0, 0, 0}},
		},
		{
			name:     "no solution",
			nums:     []int{1, 2, -2, -1},
			expected: nil,
		},
		{
			name:     "large input with duplicates",
			nums:     []int{-2, -1, 1, 1, 1, 1, 1, 2, 2, -2, -2},
			expected: [][]int{{-2, 1, 1}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := threeSum(tc.nums)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Failed %s. Got %v, expected %v", tc.name, result, tc.expected)
			}
		})
	}
}
```

main.go

```go
package main

import (
	"fmt"
	"sort"
)

func threeSum(nums []int) [][]int {
	var (
		l   = len(nums)
		ret [][]int
	)
	sort.Ints(nums)
	for i := 0; i < l-2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		k := l - 1
		for j := i + 1; j < l-1 && j < k; j++ {
			if j > i+1 && nums[j] == nums[j-1] {
				continue
			}
			for k > j+1 && nums[k] > -1*(nums[i]+nums[j]) {
				k--
			}
			if nums[i]+nums[j]+nums[k] == 0 {
				ret = append(ret, []int{nums[i], nums[j], nums[k]})
			}
		}
	}
	return ret
}
```
