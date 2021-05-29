package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	nbytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

func New(maxBytes int64, onevicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onevicted,
	}
}

func (c *Cache) Add(key string, value Value) {
	if c.cache == nil {
		c.cache = make(map[string]*list.Element)
	}
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		entry := ele.Value.(*entry)
		entry.value = value
		c.nbytes += int64(value.Len()) - int64(entry.value.Len())
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	c.nbytes += int64(len(key)) + int64(value.Len())
	if c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

func (c *Cache) Remove(key string) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.removeElement(ele)
	}
}

func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(ele *list.Element) {
	c.ll.Remove(ele)
	entry := ele.Value.(*entry)
	delete(c.cache, entry.key)
	c.nbytes -= int64(len(entry.key)) + int64(entry.value.Len())
	if c.OnEvicted != nil {
		c.OnEvicted(entry.key, entry.value)
	}
}

func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

func (c *Cache) Clear() {
	if c.OnEvicted != nil {
		for _, ele := range c.cache {
			entry := ele.Value.(*entry)
			c.OnEvicted(entry.key, entry.value)
		}
	}
	c.nbytes = 0
	c.ll = nil
	c.cache = nil
}
