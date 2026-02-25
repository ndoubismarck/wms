package memory

import (
	"fmt"
	"runtime"
	"server/internal/core/shared/types"
	"sync"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

type Options struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

type store struct {
	defaultExpiration time.Duration
	items             map[string]types.IMemoryCacheItem
	mu                sync.RWMutex
	onEvicted         func(string, interface{})
	janitor           *janitor
}

func New(ctx types.IContext) types.IMemoryCache {
	return newStoreWithJanitor(0, 0, make(map[string]types.IMemoryCacheItem))
}

func (c *store) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	c.items[k] = types.IMemoryCacheItem{
		Object:     x,
		Expiration: e,
	}
	// TODO: Calls to mu.Unlock are currently not deferred because defer adds ~200 ns (as of go1.)
	c.mu.Unlock()
}

func (c *store) SetDefault(k string, x interface{}) {
	c.Set(k, x, DefaultExpiration)
}

func (c *store) set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items[k] = types.IMemoryCacheItem{
		Object:     x,
		Expiration: e,
	}
}

func (c *store) Add(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if found {
		c.mu.Unlock()
		return fmt.Errorf("item %s already exists", k)
	}
	c.set(k, x, d)
	c.mu.Unlock()
	return nil
}

func (c *store) Replace(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if !found {
		c.mu.Unlock()
		return fmt.Errorf("item %s doesn't exist", k)
	}
	c.set(k, x, d)
	c.mu.Unlock()
	return nil
}

func (c *store) Exists(k string) bool {
	c.mu.RLock()
	_, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return false
	}
	c.mu.RUnlock()
	return true
}

func (c *store) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	// "Inlining" of get and Expired
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
	return item.Object, true
}

func (c *store) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, time.Time{}, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, time.Time{}, false
		}
		c.mu.RUnlock()
		return item.Object, time.Unix(0, item.Expiration), true
	}
	c.mu.RUnlock()
	return item.Object, time.Time{}, true
}

func (c *store) get(k string) (interface{}, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Object, true
}

func (c *store) Delete(k string) {
	c.mu.Lock()
	v, evicted := c.delete(k)
	c.mu.Unlock()
	if evicted {
		c.onEvicted(k, v)
	}
}

func (c *store) delete(k string) (interface{}, bool) {
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}
	delete(c.items, k)
	return nil, false
}

type keyAndValue struct {
	key   string
	value interface{}
}

func (c *store) DeleteExpired() {
	var evictedItems []keyAndValue
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		// "Inlining" of expired
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := c.delete(k)
			if evicted {
				evictedItems = append(evictedItems, keyAndValue{k, ov})
			}
		}
	}
	c.mu.Unlock()
	for _, v := range evictedItems {
		c.onEvicted(v.key, v.value)
	}
}

func (c *store) OnEvicted(f func(string, interface{})) {
	c.mu.Lock()
	c.onEvicted = f
	c.mu.Unlock()
}

func (c *store) Items() map[string]types.IMemoryCacheItem {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := make(map[string]types.IMemoryCacheItem, len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		// "Inlining" of Expired
		if v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		m[k] = v
	}
	return m
}

func (c *store) ItemCount() int {
	c.mu.RLock()
	n := len(c.items)
	c.mu.RUnlock()
	return n
}

func (c *store) Flush() {
	c.mu.Lock()
	c.items = map[string]types.IMemoryCacheItem{}
	c.mu.Unlock()
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *store) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(c *store) {
	c.janitor.stop <- true
}

func runJanitor(c *store, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

func newStore(de time.Duration, m map[string]types.IMemoryCacheItem) *store {
	if de == 0 {
		de = -1
	}
	c := &store{
		defaultExpiration: de,
		items:             m,
	}
	return c
}

func newStoreWithJanitor(de time.Duration, ci time.Duration, m map[string]types.IMemoryCacheItem) *store {
	c := newStore(de, m)
	if ci > 0 {
		runJanitor(c, ci)
		runtime.SetFinalizer(c, stopJanitor)
	}
	return c
}
