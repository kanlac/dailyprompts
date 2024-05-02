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

type LeveledNode struct {
	tn     *TreeNode
	father int
	level  int
}

func isCousins(root *TreeNode, x int, y int) bool {
	var (
		xNode *LeveledNode
		yNode *LeveledNode
	)

	markTarget := func(node *TreeNode, father, level int) {
		if node.Val == x {
			xNode = &LeveledNode{
				tn:     node,
				father: father,
				level:  level,
			}
		}
		if node.Val == y {
			yNode = &LeveledNode{
				tn:     node,
				father: father,
				level:  level,
			}
		}
	}

	var addChild func(*TreeNode, int, int) bool
	addChild = func(child *TreeNode, father, level int) bool {
		if xNode != nil && yNode != nil {
			return xNode.level == yNode.level && xNode.father != yNode.father
		}

		if child == nil {
			return false
		}
		markTarget(child, father, level)
		return addChild(child.Left, child.Val, level+1) || addChild(child.Right, child.Val, level+1)
	}

	markTarget(root, -1, 0)

	return addChild(root.Left, root.Val, 1) || addChild(root.Right, root.Val, 1)

}
```