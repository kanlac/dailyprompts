### 使用 Go 编写一个闭合寻址的哈希表实现

```go
package main

import (
	"fmt"
	"sync"
)

type bucket struct {
	key   string
	value string
	next  *bucket
}

type HashTable struct {
	mutex    sync.RWMutex
	capacity int
	slice    []*bucket
}

func NewHashTable(capacity int) *HashTable {
	hashTable := &HashTable{
		capacity: capacity,
		slice:    make([]*bucket, capacity),
	}
	for i := range hashTable.slice {
		hashTable.slice[i] = &bucket{}
	}
	return hashTable
}

func (h *HashTable) Get(key string) (string, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	index := h.hash(key)

	node := h.slice[index]
	for node != nil && node.key != key {
		node = node.next
	}
	if node == nil {
		return "", false
	}
	return node.value, true
}

func (h *HashTable) Set(key, value string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	index := h.hash(key)

	if h.slice[index].key == "" {
		h.slice[index].key = key
		h.slice[index].value = value
		return
	}
	if h.slice[index].key == key {
		h.slice[index].value = value
		return
	}

	node := h.slice[index]
	for node.next != nil && node.next.key != key {
		node = node.next
	}
	if node.next == nil {
		node.next = &bucket{key: key, value: value}
		return
	}

	// update
	node.next.value = value
}

func (h *HashTable) Del(key string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	index := h.hash(key)

	if h.slice[index].key == "" {
		return
	}
	if h.slice[index].key == key {
		h.slice[index] = h.slice[index].next
		return
	}

	node := h.slice[index]
	for node.next != nil && node.next.key != key {
		node = node.next
	}
	if node.next == nil {
		return
	}
	node.next = node.next.next
}

func (h *HashTable) hash(key string) int {
	var total int
	for _, r := range key {
		total += int(r)
	}
	return total % h.capacity
}

func main() {
	ht := NewHashTable(1)
	ht.Set("foo", "bar")
	ht.Set("foo0", "bar0")
	ht.Set("foo1", "bar1")

	ret, exists := ht.Get("foo")
	fmt.Printf("get foo: %s, %v\n", ret, exists)
	ret, exists = ht.Get("foo0")
	fmt.Printf("get foo0: %s, %v\n", ret, exists)
	ht.Set("foo", "foo")
	ret, exists = ht.Get("foo")
	fmt.Printf("get foo: %s, %v\n", ret, exists)
	ht.Del("foo0")
	ret, exists = ht.Get("foo0")
	fmt.Printf("get foo0: %s, %v\n", ret, exists)
	ret, exists = ht.Get("foo1")
	fmt.Printf("get foo1: %s, %v\n", ret, exists)
}


```