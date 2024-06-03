# 合并区间

https://leetcode.com/problems/merge-intervals/description/

```go
package main

import "sort"

func merge(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0] || (intervals[i][0] == intervals[j][0] && intervals[i][1] < intervals[j][1])
	})

	var ret [][]int
	for i, item := range intervals {
		last := len(ret)-1
		if i == 0 || item[0] > ret[last][1] {
			ret = append(ret, []int{item[0], item[1]})
			continue
		}
		ret[last][1] = max(item[1], ret[last][1])
	}
	return ret
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
```