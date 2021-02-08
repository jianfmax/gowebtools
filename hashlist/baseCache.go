package hashlist

import (
	"sync"
	"time"
)

type baseCache struct {
	createTime time.Time    // create 该cache的时间
	mu         sync.RWMutex // cache的全局锁
	size       int
	noLazy     bool                          // 是否是懒删除，如果为false，则只有执行Purge时才会执行，如果为true，则会固定时间内执行一次删除操作
	purgeTime  string                        // 执行删除操作的时间间隔
	maxSize    int                           // 最大的尺寸
	liveTime   *time.Duration                // 存活时间
	deleteFunc func(value interface{}) error // 执行删除操作时会执行的操作
}

// Cache Cache
type Cache interface {
	Set(key, value interface{}) error
	Get(key interface{}) (interface{}, error)
	GetThroughKeyValue(key interface{}) []KeyAndValue                           // 如果key的值相同，则存在该值,存在性能问题
	GetByFunc(key interface{}, fun func(k1, k2 interface{}) bool) []KeyAndValue // 通过一个函数比较这个值是否存在,存在性能问题
	GetAll() []KeyAndValue
	Remove(key interface{}) bool
	GetAllKey() []interface{}
	Size() int
	Has(key interface{}) bool                                                      // 是否有这个值
	HasThroughKeyValue(key interface{}) bool                                       // 如果key的值相同，则存在该值,存在性能问题
	HasByFunc(key interface{}, fun func(k1, k2 interface{}) bool) bool             // 通过一个函数比较这个值是否存在,存在性能问题
	Purge()                                                                        // 净化cache, 立刻删除过期的值
	Range(fun func(k, v interface{}, r chan interface{}) bool, r chan interface{}) // 遍历整个cache
	Init()                                                                         // 重新初始化
}
