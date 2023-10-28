### 使用 Go 编写一个闭合寻址的哈希表实现

```go
package main

import (
	"fmt"
	"sync"
)

type pair struct {
	key   string
	value string
}

type bucket []pair

type hashTable struct {
	m       sync.RWMutex
	size    int
	buckets []bucket
}

func NewHashTable(size int) *hashTable {
	buckets := make([]bucket, size)
	for i := range buckets {
		buckets[i] = make([]pair, 0)
	}
	return &hashTable{
		size:    size,
		buckets: buckets,
	}
}

func (h *hashTable) Set(key, value string) {
	h.m.Lock()
	defer h.m.Unlock()

	idx := h.hash(key)

	// 处理碰撞的情况
	for i := 0; i < len(h.buckets[idx]); i++ {
		if h.buckets[idx][i].key == key {
			h.buckets[idx][i].value = value
			return
		}
	}

	// 没有碰撞，追加
	h.buckets[idx] = append(h.buckets[idx], pair{key: key, value: value})
}

func (h *hashTable) Get(key string) (string, bool) {
	h.m.RLock()
	defer h.m.RUnlock()

	idx := h.hash(key)

	for i := 0; i < len(h.buckets[idx]); i++ {
		if h.buckets[idx][i].key == key {
			return h.buckets[idx][i].value, true
		}
	}
	return "", false
}

func (h *hashTable) Del(key string) {
	h.m.Lock()
	defer h.m.Unlock()

	idx := h.hash(key)

	for i := 0; i < len(h.buckets[idx]); i++ {
		if h.buckets[idx][i].key == key {
			h.buckets[idx] = append(h.buckets[idx][:i], h.buckets[idx][i+1:]...)
			return
		}
	}
}

func (h *hashTable) hash(key string) int {
	var total int
	for _, r := range key {
		total = total + int(r)
	}
	return total % h.size
}

func main() {
	h := NewHashTable(1)
	h.Set("foo", "bar")
	fmt.Printf("%+v\n", h)
	v, ok := h.Get("foo")
	fmt.Printf("%v, %v\n", v, ok)
	// h.Del("foo")
	// fmt.Printf("%+v\n", h)

	fmt.Println("======")

	h.Set("jack", "laurie")
	fmt.Printf("%+v\n", h)

}


```