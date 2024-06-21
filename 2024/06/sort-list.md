# 排序链表

https://leetcode.com/problems/sort-list/submissions/1295094930/

```go
package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func sortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	var (
		fast = head
		slow = head
	)
	for fast.Next != nil && fast.Next.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}
	fast = slow.Next
	slow.Next = nil

	prev := &ListNode{Val: -1}
	cur := prev
	left, right := sortList(head), sortList(fast)
	for left != nil && right != nil {
		if left.Val < right.Val {
			cur.Next = left
			left = left.Next
		} else {
			cur.Next = right
			right = right.Next
		}
		cur = cur.Next
	}
	if left != nil {
		cur.Next = left
	} else {
		cur.Next = right
	}
	return prev.Next
}
```
