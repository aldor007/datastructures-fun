package main

import "fmt"

const BUCKET_SIZE = 20

func hash(key string) (hash uint32) {
	// a jenkins one-at-a-time-hash
	// refer https://en.wikipedia.org/wiki/Jenkins_hash_function
	hash = 0
	for _, ch := range key {
		hash += uint32(ch)
		hash += hash << 10
		hash ^= hash >> 6
	}
	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15
	return
}

func getIndex(key string) (index int) {
	return int(hash(key)) % BUCKET_SIZE
}

type Node struct {
	next  *Node
	key   string
	value interface{}
}

type HashMap struct {
	values []*Node
}

func New() *HashMap {
	return &HashMap{
		values: make([]*Node, BUCKET_SIZE),
	}
}

func (h *HashMap) Get(key string) (interface{}, bool) {
	index := getIndex(key)
	if h.values[index] != nil {
		// key is on this index, but might be somewhere in linked list
		start := h.values[index]
		if start.key == key {
			return start.value, true
		}
		for start.next != nil {
			if start.key == key {
				// key matched
				return start.value, true
			}
			start = start.next
		}
	}

	// key does not exists
	return "", false
}

func (h *HashMap) Set(key string, value interface{}) {
	index := getIndex(key)

	if h.values[index] == nil {
		h.values[index] = &Node{
			next:  nil,
			key:   key,
			value: value,
		}
	} else {
		start := h.values[index]
		if start.key == key {
			start.value = value
			return
		}

		for start.next != nil {
			if start.key == key {
				start.value = value
				return
			}

			start = start.next
		}

		start.next = &Node{
			key:   key,
			value: value,
			next:  nil,
		}
	}
}

func (h *HashMap) Delete(key string) {
	index := getIndex(key)
	start := h.values[index]

	for start != nil {
		if start.key == key {
			if start.next == nil {
				h.values[index] = nil
				return
			} else {
				h.values[index] = start.next
			}
		}
		start = start.next
	}
}

func main() {
	h := New()
	h.Set("key-1", "value")
	fmt.Println(h.Get("key-1"))
	h.Delete("key-1")
	fmt.Println(h.Get("key-1"))
}