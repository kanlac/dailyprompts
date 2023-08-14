# 找到二叉树的左侧兄弟节点

```go
package main

import "fmt"

type TreeNode struct {
	Val        int
	LeftChild  *TreeNode
	RightChild *TreeNode

	Cousin *TreeNode
}

type BFSNode struct {
	node  *TreeNode
	depth int
}

func BuildTree(root *TreeNode) *TreeNode {
	if root == nil {
		return root
	}

	queue := []*BFSNode{{node: root, depth: 0}}
	lastNode := &BFSNode{depth: -1}

	for len(queue) > 0 {
		head := queue[0]
		dep := head.depth
		queue = queue[1:]

		if head.node.LeftChild != nil {
			queue = append(queue, &BFSNode{node: head.node.LeftChild, depth: dep + 1})
		}
		if head.node.RightChild != nil {
			queue = append(queue, &BFSNode{node: head.node.RightChild, depth: dep + 1})
		}

		if lastNode.depth == dep {
			head.node.Cousin = lastNode.node
		}
		lastNode = head
	}

	return root
}

func main() {
	four := &TreeNode{Val: 4}
	five := &TreeNode{Val: 5}
	six := &TreeNode{Val: 6}
	two := &TreeNode{Val: 2, RightChild: four}
	three := &TreeNode{Val: 3, LeftChild: five, RightChild: six}
	one := &TreeNode{Val: 1, LeftChild: two, RightChild: three}
	r := BuildTree(one)
	fmt.Println(r)
}


```