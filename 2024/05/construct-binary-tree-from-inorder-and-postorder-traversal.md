# 从中序与后序遍历序列构造二叉树

https://leetcode.com/problems/construct-binary-tree-from-inorder-and-postorder-traversal

```go
func buildTree(inorder []int, postorder []int) *TreeNode {
	l := len(inorder)
	if l == 0 {
		return nil
	}
	rootVal := postorder[l-1]

	var rootIdx int
	for i, n := range inorder {
		if n == rootVal {
			rootIdx = i
			break
		}
	}
	leftNum := rootIdx

	return &TreeNode{
		Val:   rootVal,
		Left:  buildTree(inorder[:rootIdx], postorder[:leftNum]),
		Right: buildTree(inorder[rootIdx+1:], postorder[leftNum:l-1]),
	}
}
```