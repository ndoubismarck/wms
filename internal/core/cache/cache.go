package cache

import (
	"server/internal/core/cache/memory"
	"server/internal/core/shared/types"
)

type cache struct {
	ctx    types.IContext
	memory types.IMemoryCache
}

func New(ctx types.IContext) types.ICache {
	return &cache{
		ctx:    ctx,
		memory: memory.New(ctx),
	}
}

func (k *cache) Memory(name string) types.IMemoryCache {
	val, exists := k.memory.Get(name)
	if !exists {
		obj := memory.New(k.ctx)
		k.memory.Set(name, obj, 0)
		return obj
	}
	obj, valid := val.(types.IMemoryCache)
	if !valid {
		k.ctx.Logger().Fatal("invalid memory cache value")
	}
	return obj
}
