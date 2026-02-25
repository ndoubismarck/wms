package types

import "time"

type ICache interface {
	Memory(name string) IMemoryCache
}

type IMemoryCacheItem struct {
	Object     interface{} `json:"object"`
	Expiration int64       `json:"expiration"`
}

type IMemoryCache interface {
	Set(k string, x interface{}, d time.Duration)

	SetDefault(k string, x interface{})

	Add(k string, x interface{}, d time.Duration) error

	Replace(k string, x interface{}, d time.Duration) error

	Exists(k string) bool

	Get(k string) (interface{}, bool)

	GetWithExpiration(k string) (interface{}, time.Time, bool)

	Delete(k string)

	DeleteExpired()

	OnEvicted(f func(string, interface{}))
}

func (item IMemoryCacheItem) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}
