# 跳跃游戏

https://leetcode.com/problems/jump-game/description/

```go
package main

func canJump(nums []int) bool {
	var (
		l         = len(nums)
		rightmost int
	)
	for i := 0; i < min(l, rightmost+1); i++ {
		rightmost = max(rightmost, i+nums[i])
		if rightmost >= l-1 {
			return true
		}
	}
	return false
}
```