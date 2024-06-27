# 二叉搜索树中第K小的元素

https://leetcode.com/problems/kth-smallest-element-in-a-bst/description/

main_test.go

```go
package main

import (
    "testing"
)

func TestKthSmallest(t *testing.T) {
    tests := []struct {
        root *TreeNode
        k    int
        want int
    }{
        {
            root: &TreeNode{3, &TreeNode{1, nil, &TreeNode{2, nil, nil}}, &TreeNode{4, nil, nil}},
            k:    1,
            want: 1,
        },
        {
            root: &TreeNode{5, &TreeNode{3, &TreeNode{2, &TreeNode{1, nil, nil}, nil}, &TreeNode{4, nil, nil}}, &TreeNode{6, nil, nil}},
            k:    3,
            want: 3,
        },
        {
            root: &TreeNode{2, &TreeNode{1, nil, nil}, &TreeNode{3, nil, nil}},
            k:    2,
            want: 2,
        },
        {
            root: &TreeNode{1, nil, &TreeNode{2, nil, nil}},
            k:    2,
            want: 2,
        },
    }

    for _, tt := range tests {
        t.Run("", func(t *testing.T) {
            got := kthSmallest(tt.root, tt.k)
            if got != tt.want {
                t.Errorf("kthSmallest() = %v, want %v", got, tt.want)
            }
        })
    }
}

```
    
solution
    
```go
func kthSmallest(root *TreeNode, k int) int {
    var stack []*TreeNode
    for {
        for root != nil {
            stack = append(stack, root)
            root = root.Left
        }   
        root, stack, k = stack[len(stack)-1], stack[:len(stack)-1], k-1
        if k == 0 {
            return root.Val
        }
        root = root.Right
    }
    return -1
}
```