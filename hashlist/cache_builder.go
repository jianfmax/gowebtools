package hashlist

import (
	"sync"
	"time"
)

type CacheType string

const (
	TYPE_HASHLIST = "hash_list"
)

// CacheBuilder 新建缓存的构建器
type CacheBuilder struct {
	createTime time.Time
	tp         CacheType
	maxSize    int
	liveTime   *time.Duration // 一个节点的存活时间
	mu         sync.RWMutex
	noLazy     bool                          // 是否是懒删除，如果为false，则只有执行Purge时才会执行，如果为true，则会固定时间内执行一次删除操作
	deleteFunc func(value interface{}) error // 执行删除操作时会执行的操作
	purgeTime  string                        // 执行删除操作的时间间隔
}

// New 创建一个构造器
func New() *CacheBuilder {
	return &CacheBuilder{
		createTime: time.Now(),
		tp:         TYPE_HASHLIST,
	}
}

// MaxSize 设置最大值
func (cb *CacheBuilder) MaxSize(size int) *CacheBuilder {
	cb.maxSize = size
	return cb
}

// SetLiveTime 一个节点的存活时间
func (cb *CacheBuilder) SetLiveTime(liveTime *time.Duration) *CacheBuilder {
	cb.liveTime = liveTime
	return cb
}

// SetLiveTime 一个节点的存活时间
func (cb *CacheBuilder) AddDeleteFunc(f func(value interface{}) error) *CacheBuilder {
	cb.deleteFunc = f
	return cb
}

// PurgeOnTime 定期删除过期节点
func (cb *CacheBuilder) PurgeOnTime(purgeTime string) *CacheBuilder {
	cb.noLazy = true
	cb.purgeTime = purgeTime
	return cb
}

// evictType 设置为指定的类型
func (cb *CacheBuilder) evictType(tp CacheType) *CacheBuilder {
	cb.tp = tp
	return cb
}

// HashList 设置类型为HashList
func (cb *CacheBuilder) HashList() *CacheBuilder {
	return cb.evictType(TYPE_HASHLIST)
}

// New 创建一个指定类型的cache
func (cb *CacheBuilder) New() Cache {
	switch cb.tp {
	case TYPE_HASHLIST:
		return newHashList(cb)
	default:
		return newHashList(cb)
	}
}

func buildCache(c *baseCache, cb *CacheBuilder) {
	c.maxSize = cb.maxSize
	c.noLazy = cb.noLazy
	c.purgeTime = cb.purgeTime
	c.liveTime = cb.liveTime
	c.deleteFunc = cb.deleteFunc
}
