# 二叉树的堂兄弟节点

https://leetcode.com/problems/cousins-in-binary-tree/description/

main_test.go
```go
package main

import "testing"

func TestIsCousins(t *testing.T) {
	testCases := []struct {
		name string
		root *TreeNode
		x    int
		y    int
		want bool
	}{
		{
			name: "case1",
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val:  2,
					Left: &TreeNode{Val: 4},
				},
				Right: &TreeNode{
					Val:   3,
					Right: &TreeNode{Val: 5},
				},
			},
			x:    4,
			y:    5,
			want: true,
		},
		{
			name: "case2",
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val:  2,
					Left: &TreeNode{Val: 4},
				},
				Right: &TreeNode{
					Val: 3,
				},
			},
			x:    4,
			y:    3,
			want: false,
		},
		{
			name: "case2",
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val: 2,
					Right: &TreeNode{
						Val: 4,
					},
				},
				Right: &TreeNode{
					Val: 3,
					Right: &TreeNode{
						Val: 5,
					},
				},
			},
			x:    5,
			y:    4,
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isCousins(tc.root, tc.x, tc.y)
			if got != tc.want {
				t.Errorf("TestIsCousins(%s) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}
```

main.go
```go
package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func isCousins(root *TreeNode, x int, y int) bool {
	var (
		// 找到的第一个节点的父亲和深度
		fatherOfFirstOccurrence = -1
		levelOfFirstOccurrence  = -1
		dfs                     func(*TreeNode, int, int) bool
	)

	// 若 match，且找到了第一个，直接返回结果
	// 若 match，且未找到第一个，标记并继续 dfs
	// 若当前节点 not match，继续 dfs
	dfs = func(node *TreeNode, fatherOfThis, levelOfThis int) bool {
		if node == nil {
			return false
		}

		if node.Val == x || node.Val == y {
			if fatherOfFirstOccurrence >= 0 {
				return fatherOfThis != fatherOfFirstOccurrence && levelOfThis == levelOfFirstOccurrence
			}
			fatherOfFirstOccurrence = fatherOfThis
			levelOfFirstOccurrence = levelOfThis
		}
		return dfs(node.Left, node.Val, levelOfThis+1) || dfs(node.Right, node.Val, levelOfThis+1)
	}

	return dfs(root, -1, 0)
}

```