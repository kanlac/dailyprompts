# 最大子数组的和

```go
package main

import (
	"math"
)

func maxSubArray(nums []int) int {
	if len(nums) == 0 {
		return math.MinInt
	}

	var (
		l          = len(nums)
		mid        = l / 2
		maxToLeft  int
		maxToRight int
	)

	for i, leftSum := mid-1, 0; i >= 0; i-- {
		leftSum += nums[i]
		maxToLeft = max(maxToLeft, leftSum)
	}
	for i, rightSum := mid+1, 0; i < l; i++ {
		rightSum += nums[i]
		maxToRight = max(maxToRight, rightSum)
	}
	maxAcross := maxToLeft + nums[mid] + maxToRight

	return maxInThree(maxAcross, maxSubArray(nums[:mid]), maxSubArray(nums[mid+1:]))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxInThree(a, b, c int) int {
	return max(a, max(b, c))
}
```