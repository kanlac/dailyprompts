# 列表转树

```golang
package main

import "fmt"

type treeNode struct {
	id     int
	parent *treeNode
	value  int
}

func generateInput() []*treeNode {
    a := &treeNode{
        id: 0,
        parent: nil,
        value: 100,
    }
    b := &treeNode{
        id: 1,
        parent: a,
        value: 101,
    }
    c := &treeNode{
        id: 2,
        parent: a,
        value: 102,
    }
	return []*treeNode{b, c, a}
}

func buildTree(list []*treeNode) map[int]*treeNode {
	m := make(map[int]*treeNode)
	// build map
	for _, ln := range list {
		m[ln.id] = &treeNode{id: ln.id, value: ln.value}
	}
	// build tree
	for _, ln := range list {
		if ln.parent == nil {
			continue
		}
		m[ln.id].parent = m[ln.parent.id]
	}
	return m
}

func main() {
	root := buildTree(generateInput())
	fmt.Printf("%+v", root)
}

```