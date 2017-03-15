package orderedmap

import (
    "encoding/json"
    "fmt"
    "hash/fnv"
)

const N_BUCKETS = 10

type node struct {
    key   string
    value interface{}

    prev  *node
    next  *node
}

type OrderedMap struct {
    size    uint32
    buckets [][]*node

    first   *node
    last    *node
}

func NewOrderedMap() *OrderedMap {
    return &OrderedMap{
        size: 0,
        buckets: make([][]*node, N_BUCKETS),
        first: nil,
        last: nil,
    }
}

func (m *OrderedMap) Get(key string) (interface{}, bool) {
    node, exists := m.getnode(key)
    if exists {
        return node.value, true
    } else {
        return nil, false
    }
}

func (m *OrderedMap) GetIndex(i uint32) (interface{}, bool) {
    node, exists := m.getnodeat(i)
    if exists {
        return node.value, true
    } else {
        return nil, false
    }
}

func (m *OrderedMap) Set(key string, value interface{}) {
    n, exists := m.getnode(key)
    if exists {
        n.value = value
        return
    }
    n = &node{
        key: key,
        value: value,

        prev: m.last,
        next: nil,
    }
    if m.last != nil {
        m.last.next = n
    } else {
        m.first = n
    }
    m.last = n
    m.size++
    m.rehash()
}

func (m *OrderedMap) Remove(key string) (interface{}, bool) {
    n, exists := m.getnode(key)
    if !exists {
        return nil, false
    }
    if n.prev != nil {
        n.prev.next = n.next
        n.prev = nil
    } else {
        m.first = n.next
    }
    if n.next != nil {
        n.next.prev = n.prev
        n.next = nil
    } else {
        m.last = n.prev
    }
    m.size--
    m.rehash()
    return n.value, true
}

func (m *OrderedMap) RemoveIndex(i uint32) (interface{}, bool) {
    n, exists := m.getnodeat(i)
    if !exists {
        return nil, false
    }
    if n.prev != nil {
        n.prev.next = n.next
        n.prev = nil
    } else {
        m.first = n.next
    }
    if n.next != nil {
        n.next.prev = n.prev
        n.next = nil
    } else {
        m.last = n.prev
    }
    m.size--
    m.rehash()
    return n.value, true
}

func (m *OrderedMap) HasKey(key string) bool {
    _, has := m.Get(key)
    return has
}

func (m *OrderedMap) HasValue(value interface{}) bool {
    for _, bucket := range m.buckets {
        for _, node := range bucket {
            if node.value == value || (node.value == nil && value == nil) {
                return true
            }
        }
    }
    return false
}

func (m *OrderedMap) Size() uint32 {
    return m.size
}

func (m *OrderedMap) Each(f func(key string, value interface{})) {
    for node := m.first; node != nil; node = node.next {
        f(node.key, node.value)
    }
}

func (m *OrderedMap) Print() {
    fmt.Printf("OrderedMap(#elems = %d) {\n", m.size)
    for n := m.first; n != nil; n = n.next {
        fmt.Printf("  %s = %s\n", n.key, n.value)
    }
    fmt.Println("}")
}

func (m *OrderedMap) MarshalJSON() ([]byte, error) {
    data := make([]byte, 0)
    data = append(data, '{')
    for n := m.first; n != nil; n = n.next {
        k, _ := json.Marshal(n.key)
        v, _ := json.Marshal(n.value)
        data = append(data, k...)
        data = append(data, ": "...)
        data = append(data, v...)
        if n.next != nil {
            data = append(data, ',')
        }
    }
    return append(data, '}'), nil
}

func (m *OrderedMap) getnode(key string) (*node, bool) {
    if m.size == 0 {
        return nil, false
    }
    h := m.hash(key)
    bucket := m.buckets[h]
    for i := range bucket {
        if bucket[i].key == key {
            return bucket[i], true
        }
    }
    return nil, false
}

func (m *OrderedMap) getnodeat(i uint32) (*node, bool) {
    if m.size < i {
        return nil, false
    }
    for n := m.first; n != nil; n = n.next {
        if i == 0 {
            return n, true
        }
        i--
    }
    return nil, false
}

func (m *OrderedMap) rehash() {
    m.buckets = make([][]*node, m.size)
    for n := m.first; n != nil; n = n.next {
        h := m.hash(n.key)
        m.buckets[h] = append(m.buckets[h], n)
    }
}

func (m *OrderedMap) hash(s string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(s))
    return h.Sum32() % m.size
}