# 最大子数组的和

```go
package main

import (
	"fmt"
	"math"
)

func maxSubArray(nums []int) int {
	l := len(nums)
	if l == 0 {
		return math.MinInt
	}

	mid, lMaxSum, rMaxSum := l/2, 0, 0
	for i, curSum := mid-1, 0; i >= 0; i-- {
		curSum += nums[i]
		lMaxSum = maxInTwo(lMaxSum, curSum)
	}
	for i, curSum := mid+1, 0; i < l; i++ {
		curSum += nums[i]
		rMaxSum = maxInTwo(rMaxSum, curSum)
	}
	maxCrossingSum := nums[mid] + lMaxSum + rMaxSum

	return maxInThree(maxSubArray(nums[:mid]), maxSubArray(nums[mid+1:]), maxCrossingSum)
}

func maxInTwo(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxInThree(a, b, c int) int {
	return maxInTwo(maxInTwo(a, b), c)
}

func main() {
	fmt.Println(maxSubArray([]int{-2, 1, -3, 4, -1, 2, 1, -5, 4})) // 6
	fmt.Println(maxSubArray([]int{1}))                             // 1
	fmt.Println(maxSubArray([]int{5, 4, -1, 7, 8}))                // 23
	fmt.Println(maxSubArray([]int{-2, 1}))                         // 1
	fmt.Println(maxSubArray([]int{-1, -2}))                        // -1
}
```