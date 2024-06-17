# 最近公共祖先

剑指 Offer 68，[LeetCode 236](https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-tree/description/)，[LeetCode CN](https://leetcode.cn/problems/er-cha-shu-de-zui-jin-gong-gong-zu-xian-lcof/description/)。

```go
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	if root.Val == p.Val || root.Val == q.Val {
		return root
	}
	left := lowestCommonAncestor(root.Left, p, q)
	right := lowestCommonAncestor(root.Right, p, q)
	if right == nil {
		return left
	}
	if left == nil {
		return right
	}
	return root
}
```