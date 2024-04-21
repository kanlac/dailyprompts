# 找到二叉树的左侧兄弟节点

```go
package main

type TreeNode struct {
	Val        int
	LeftChild  *TreeNode
	RightChild *TreeNode

	Cousin *TreeNode
}

type leveledNode struct {
	tn    *TreeNode
	level int
}

func BuildTree(root *TreeNode) *TreeNode {
	if root == nil {
		return root
	}
	
	queue := []leveledNode{{tn: root, level: 0}}
	var lastPoped leveledNode
	for len(queue) > 0 {
		pop := queue[0]
		queue = queue[1:]

		if lastPoped.tn != nil && lastPoped.level == pop.level {
			pop.tn.Cousin = lastPoped.tn
		}

		if pop.tn.LeftChild != nil {
			queue = append(queue, leveledNode{tn: pop.tn.LeftChild, level: pop.level + 1})
		}
		if pop.tn.RightChild != nil {
			queue = append(queue, leveledNode{tn: pop.tn.RightChild, level: pop.level + 1})
		}

		lastPoped = pop
	}
	return root
}
```
