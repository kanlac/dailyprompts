# 找到二叉树的左侧兄弟节点

main_test.go

```go
package main

import (
	"testing"
)

func TestBuildTree(t *testing.T) {
	tests := []struct {
		name        string
		root        *TreeNode
		wantCousins map[int]int // key为节点值，value为该节点预期的Cousin节点的ID
	}{
		{
			name:        "single node",
			root:        &TreeNode{Val: 1},
			wantCousins: map[int]int{},
		},
		{
			name: "complete binary tree",
			root: &TreeNode{
				Val:        1,
				LeftChild:  &TreeNode{Val: 2},
				RightChild: &TreeNode{Val: 3},
			},
			wantCousins: map[int]int{
				3: 2,
			},
		},
		{
			name: "complete binary tree 2",
			root: &TreeNode{
				Val: 1,
				LeftChild: &TreeNode{
					Val:        2,
					RightChild: &TreeNode{Val: 4},
				},
				RightChild: &TreeNode{
					Val: 3,
					LeftChild: &TreeNode{
						Val: 5,
					},
					RightChild: &TreeNode{
						Val: 6,
					},
				},
			},
			wantCousins: map[int]int{
				3: 2,
				5: 4,
				6: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BuildTree(tt.root)

			queue := []*TreeNode{tt.root}
			if tt.root.LeftChild != nil {
				queue = append(queue, tt.root.LeftChild)
			}
			if tt.root.RightChild != nil {
				queue = append(queue, tt.root.RightChild)
			}

			for len(queue) > 0 {
				popped := queue[0]
				queue = queue[1:]

				if _, ok := tt.wantCousins[popped.Val]; ok {
					if popped.Cousin == nil || popped.Cousin.Val != tt.wantCousins[popped.Val] {
						t.Errorf("node %+v wanted cousin %d", popped, tt.wantCousins[popped.Val])
					}
				}
			}
		})
	}
}

```

main.go

```go
package main

type TreeNode struct {
	Val        int
	LeftChild  *TreeNode
	RightChild *TreeNode

	Cousin *TreeNode
}

type levelNode struct {
	*TreeNode
	level int
}

func BuildTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}

	var (
		queue     = []*levelNode{{TreeNode: root, level: 0}}
		lastLevel = -1
		lastNode  *TreeNode
	)

	for len(queue) > 0 {
		pop := queue[0]
		queue = queue[1:]

		if pop.level == lastLevel {
			pop.Cousin = lastNode
		}
		lastLevel, lastNode = pop.level, pop.TreeNode

		if pop.LeftChild != nil {
			queue = append(queue, &levelNode{TreeNode: pop.LeftChild, level: pop.level + 1})
		}
		if pop.RightChild != nil {
			queue = append(queue, &levelNode{TreeNode: pop.RightChild, level: pop.level + 1})
		}
	}
	return root
}
```
