# 列表转树

## Q

有一个节点的列表，它们都是按层级顺序排好的，每个节点有一个 parentID。请从这个列表中构建一颗树。

## 测试用例及其作答模版

main.go:

```go
package main

type ListNode struct {
 id       int
 parentID int
 val      string
}

type TreeNode struct {
 id       int
 children []*TreeNode
 val      string
}

func BuildTree(list []ListNode) *TreeNode {
 return &TreeNode{} // Your implementation goes here
}
```

main_test.go:

```go
package main

import (
 "reflect"
 "testing"
)

func TestBuildTree(t *testing.T) {
 testCases := []struct {
  input    []ListNode
  expected *TreeNode
 }{
  {
   input: []ListNode{
    {id: 1, parentID: 0, val: "root"},
    {id: 2, parentID: 1, val: "child1"},
    {id: 3, parentID: 1, val: "child2"},
   },
   expected: &TreeNode{
    id:  1,
    val: "root",
    children: []*TreeNode{
     {id: 2, val: "child1", children: []*TreeNode{}},
     {id: 3, val: "child2", children: []*TreeNode{}},
    },
   },
  },
  {
   input: []ListNode{
    {id: 1, parentID: 0, val: "root"},
   },
   expected: &TreeNode{
    id:       1,
    val:      "root",
    children: []*TreeNode{},
   },
  },
  {
   input: []ListNode{
    {id: 1, parentID: 0, val: "root"},
    {id: 2, parentID: 1, val: "child1"},
    {id: 3, parentID: 2, val: "grandchild1"},
   },
   expected: &TreeNode{
    id:  1,
    val: "root",
    children: []*TreeNode{
     {
      id:  2,
      val: "child1",
      children: []*TreeNode{
       {id: 3, val: "grandchild1", children: []*TreeNode{}},
      },
     },
    },
   },
  },
 }

 for i, tc := range testCases {
  result := BuildTree(tc.input)
  if !reflect.DeepEqual(result, tc.expected) {
   t.Errorf("Test case %d failed; got=%#v, expected=%#v", i, result, tc.expected)
  }
 }
}
```

## 实现

```go
func BuildTree(list []ListNode) *TreeNode {
	l := len(list)
	if l < 1 {
		return nil
	}

	root := &TreeNode{
		id:  list[0].id,
		val: list[0].val,
	}
	
	m := make(map[int]*TreeNode)
	m[root.id] = root

	for i := 1; i < l; i++ {
		tn := &TreeNode{
			id:  list[i].id,
			val: list[i].val,
		}
		m[list[i].id] = tn

		if list[i].parentID > 0 { // 认为 id 一定大于 0
			parent := m[list[i].parentID]
			if parent.children == nil {
				parent.children = make([]*TreeNode, 0)
			}
			parent.children = append(parent.children, tn)
		}
	}

	return root
}
```
