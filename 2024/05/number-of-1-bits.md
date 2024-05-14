# 汉明重量

https://leetcode.com/problems/number-of-1-bits/

```go
func hammingWeight(n int) int {
	var weight int
	for ; n > 0; weight++ {
		n &= n-1
	}
	return weight
}
```
