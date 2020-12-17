package caches

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	data map[string]*value
	options Options
	status *Status
	lock *sync.RWMutex
}

func NewCache() *Cache {
	return NewCacheWith(DefaultOptions())
}

func NewCacheWith(options Options) *Cache {
	return &Cache{
		data: make(map[string]*value, 256),
		options: options,
		status: newStatus(),
		lock: &sync.RWMutex{},
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return nil, false
	}

	if !value.alive() {
		c.lock.RUnlock()
		c.Delete(key)
		c.lock.RLock()
		return nil, false
	}

	return value.visit(), true
}

func (c *Cache) Set(key string, value []byte) error {
	return c.SetWithTTL(key, value, NeverDie)
}

func (c *Cache) SetWithTTL(key string, value []byte, ttl int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if oldValue, ok := c.data[key]; ok {
		c.status.subEntry(key, oldValue.data)
	}

	if !c.checkEntrySize(key, value) {
		if oldValue, ok := c.data[key]; ok {
			c.status.addEntry(key, oldValue.data)
		}

		return errors.New("the entry size will exceed if you set thie entry")
	}

	c.status.addEntry(key, value)
	c.data[key] = newValue(value, ttl)
	return nil
}

func (c *Cache) checkEntrySize(newKey string, newValue []byte) bool {
	return c.status.entrySize()+int64(len(newKey))+int64(len(newValue)) <= c.options.MaxEntrySize*1024*1024
}

func (c *Cache) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if oldValue, ok := c.data[key]; ok {
		c.status.subEntry(key, oldValue.data)
		delete(c.data, key)
	}
}

func (c *Cache) Status() Status {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return *c.status
}

func (c *Cache) gc() {
	c.lock.Lock()
	defer c.lock.Unlock()

	count := 0
	for key, value := range c.data {
		if !value.alive() {
			c.status.subEntry(key, value.data)
			delete(c.data, key)
			count++
			if count >= c.options.MaxGcCount {
				break
			}
		}
	}
}

func (c *Cache) AutoGc() {
	go func() {
		ticker := time.NewTicker(time.Duration(c.options.GcDuration)*time.Minute)
		for {
			select {
			case <-ticker.C:
				c.gc()
			}
		}
	}()
}