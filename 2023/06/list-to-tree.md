# 列表转树

### Q
有一个节点的列表，它们都是按层级顺序排好的，每个节点有一个 parentID。请从这个列表中构建一颗树。

### A

```golang
package main

import "fmt"

type ListNode struct {
	id       int
	val      int
	parentID int
}

type TreeNode struct {
	id     int
	val    int
	parent *TreeNode
}

func buildTree(list []ListNode) map[int]*TreeNode {
	m := make(map[int]*TreeNode)
	for _, n := range list {
		parent := m[n.parentID]
		this := &TreeNode{id: n.id, val: n.val, parent: parent}
		m[n.id] = this
	}
	return m
}

func main() {
	list := []ListNode{
		{
			id:       0,
			val:      3,
			parentID: -1,
		},
		{
			id:       1,
			val:      2,
			parentID: 0,
		},
		{
			id:       2,
			val:      4,
			parentID: 0,
		},
		{
			id:       3,
			val:      1,
			parentID: 1,
		},
		{
			id:       4,
			val:      5,
			parentID: 1,
		},
	}
	tree := buildTree(list)
	fmt.Println(tree)
}

```