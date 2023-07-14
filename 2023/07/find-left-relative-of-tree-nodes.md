# 找到二叉树的左侧兄弟节点

```go
package main

type SrcTreeNode struct {
	val      int
	leftChi  *SrcTreeNode
	rightChi *SrcTreeNode
}

type TrgTreeNode struct {
	srcTreeNode *SrcTreeNode
	level    int
	leftRel  *TrgTreeNode
}

func BuildTree(root *SrcTreeNode) *TrgTreeNode {
	if root == nil {
		return nil
	}

	trgRoot := &TrgTreeNode{
		srcTreeNode: root,
		level: 1,
	}
	stack := []*TrgTreeNode{trgRoot}

	var lastNode *TrgTreeNode
	for len(stack) != 0 { // BFS
		head := stack[0]
		stack = stack[1:]
		if lastNode != nil && head.level == lastNode.level {
			head.leftRel = lastNode
		}

		// enqueue
		if head.srcTreeNode.leftChi != nil {
			stack = append(stack, &TrgTreeNode{
				srcTreeNode: head.srcTreeNode.leftChi,
				level: head.level + 1,
			})
		}
		if head.srcTreeNode.rightChi != nil {
			stack = append(stack, &TrgTreeNode{
				srcTreeNode: head.srcTreeNode.rightChi,
				level: head.level + 1,
			})
		}
		lastNode = head
	}
	return trgRoot
}

```