package lfu

import (
	"container/heap"

	"github.com/lixvyang/dutcache/pkg/cachestruct"
)

// Cache is an LFU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64
	nbytes   int64
	cache    map[string]*entry
	pq       priorityQueue
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value cache.Value)
}

type entry struct {
	key      string
	value    cache.Value
	priority int // LFU priority
	index    int // index in the priority queue
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, cache.Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		cache:     make(map[string]*entry),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value cache.Value, ok bool) {
	if ent, ok := c.cache[key]; ok {
		value = ent.value
		c.incrementPriority(ent)
		return value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	if len(c.cache) == 0 {
		return
	}
	ent := heap.Pop(&c.pq).(*entry)
	delete(c.cache, ent.key)
	c.nbytes -= int64(len(ent.key)) + int64(ent.value.Len())
	if c.OnEvicted != nil {
		c.OnEvicted(ent.key, ent.value)
	}
}

func (c *Cache) Add(key string, value cache.Value) {
	if ent, ok := c.cache[key]; ok {
		c.nbytes += int64(value.Len()) - int64(ent.value.Len())
		ent.value = value
		c.incrementPriority(ent)
	} else {
		ent := &entry{
			key:      key,
			value:    value,
			priority: 1,
		}
		c.cache[key] = ent
		heap.Push(&c.pq, ent)
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return len(c.cache)
}

func (c *Cache) incrementPriority(ent *entry) {
	ent.priority++
	heap.Fix(&c.pq, ent.index)
}

// priorityQueue implements heap.Interface and holds cache entries.
type priorityQueue []*entry

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the entry with the lowest priority
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	ent := x.(*entry)
	ent.index = n
	*pq = append(*pq, ent)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	ent := old[n-1]
	ent.index = -1 // for safety
	*pq = old[0 : n-1]
	return ent
}
