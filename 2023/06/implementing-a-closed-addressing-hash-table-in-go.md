### 使用 Go 编写一个闭合寻址的哈希表实现

```go
package main

import "fmt"

const ArraySize = 7

// HashTable struct
type HashTable struct {
	array [ArraySize]*bucket
}

// bucket struct
type bucket struct {
	key   string
	value int
	next  *bucket
}

// Insert will take in a key and value
// hash the key
// add the key and value to the array at hashed index
func (h *HashTable) Insert(key string, value int) {
	index := hash(key)
	h.array[index].insert(key, value)
}

// Search will take in a key and return value
// hash the key
// search the array at hashed index
func (h *HashTable) Search(key string) (int, bool) {
	index := hash(key)
	return h.array[index].search(key)
}

// hash
func hash(key string) int {
	sum := 0
	for _, v := range key {
		sum += int(v)
	}
	return sum % ArraySize
}

// insert
func (b *bucket) insert(k string, v int) {
	if b.key == "" {
		b.key = k
		b.value = v
		return
	}

	ptr := b
	for ptr.next != nil {
		ptr = ptr.next
	}
	ptr.next = &bucket{key: k, value: v}
}

// search
func (b *bucket) search(k string) (int, bool) {
	ptr := b
	while(ptr != nil) {
		if ptr.key == k {
			return ptr.value, true
		}
		ptr = ptr.next
	}
	return 0, false
}

// Init initializes the hashtable
func Init() *HashTable {
	result := &HashTable{}
	for i := range result.array {
		result.array[i] = &bucket{}
	}
	return result
}

func main() {
	h := Init()
	h.Insert("ERIC", 589)
	h.Insert("KEN", 522)
	h.Insert("RON", 101)
	h.Insert("SAM", 229)
	h.Insert("SUN", 943)
	fmt.Println(h.Search("KEN"))
	fmt.Println(h.Search("ERIC"))
	fmt.Println(h.Search("SUN"))
	fmt.Println(h.Search("MIKE"))
}

```