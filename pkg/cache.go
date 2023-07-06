package dutcache

import (
	"sync"

	c "github.com/lixvyang/dutcache/pkg/cache"
	"github.com/lixvyang/dutcache/pkg/cache/lfu"
)

type cache struct {
	mu         sync.Mutex
	dataStruct c.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.dataStruct == nil {
		c.dataStruct = lfu.New(c.cacheBytes, nil)
	}
	c.dataStruct.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.dataStruct == nil {
		return
	}

	if v, ok := c.dataStruct.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
