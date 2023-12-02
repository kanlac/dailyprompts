# 分治法

[108) Convert Sorted Array to Binary Search Tree](https://leetcode.com/problems/convert-sorted-array-to-binary-search-tree/description/)

```go
package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func sortedArrayToBST(nums []int) *TreeNode {
	l := len(nums)
	if l == 0 {
		return nil
	}

	mid := l / 2
	return &TreeNode{
		Val:   nums[mid],
		Left:  sortedArrayToBST(nums[:mid]),
		Right: sortedArrayToBST(nums[mid+1:]),
	}
}
```

[53) Maximum Subarray](https://leetcode.com/problems/maximum-subarray/description/)

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

[106) 106. Construct Binary Tree from Inorder and Postorder Traversal](https://leetcode.com/problems/construct-binary-tree-from-inorder-and-postorder-traversal/)

```go
package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func buildTree(inorder []int, postorder []int) *TreeNode {
	m := make(map[int]int)
	for i := range inorder {
		m[inorder[i]] = i
	}

	var construct func(inStart, inEnd, postStart, postEnd int) *TreeNode
	construct = func(inStart, inEnd, postStart, postEnd int) *TreeNode {
		if inStart > inEnd || postStart > postEnd {
			return nil
		}

		root := &TreeNode{Val: postorder[postEnd]}
		leftNum := m[root.Val] - inStart

		root.Left = construct(inStart, m[root.Val]-1, postStart, postStart+leftNum-1)
		root.Right = construct(m[root.Val]+1, inEnd, postStart+leftNum, postEnd-1)
		return root
	}

	return construct(0, len(inorder)-1, 0, len(postorder)-1)
}
```
