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

import "sort"

func threeSum(nums []int) [][]int {
	sort.Ints(nums)

	var (
		ret [][]int
		l   = len(nums)
	)

	for i := 0; i < l; i++ {
		if i > 0 && nums[i-1] == nums[i] {
			continue
		}
		for j, k := i+1, l-1; j < l; j++ {
			if j > i+1 && nums[j-1] == nums[j] {
				continue
			}
			// 剪枝：确定 i, j 后，只用找到第一个 k 即可
			// 剪枝：在 i 确定时，随着 j 迭代，k 的右边界也跟着向前推进
			for k > j && nums[k] > -1*(nums[i]+nums[j]) {
				k--
			}
			if k <= j {
				break
			}
			if nums[k] == -1*(nums[i]+nums[j]) {
				ret = append(ret, []int{nums[i], nums[j], nums[k]})
			}
		}
	}
	return ret
}
```