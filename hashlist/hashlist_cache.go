package hashlist

import (
	"container/list"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"time"

	"github.com/robfig/cron/v3"
)

// HashListCache hash list 组成的缓存
type HashListCache struct {
	baseCache
	items map[interface{}]*list.Element
	lists *list.List
}
type simpleItem struct {
	Key        interface{}
	CreateTime *time.Time // 创建该节点的时间
	Value      interface{}
	DeleteTime *time.Time             // 被删除的时间
	Msg        map[string]interface{} // 一些其余的信息
}

type KeyAndValue struct {
	Key   interface{}
	Value interface{}
}

// new 创建一个新的hashList
func newHashList(cb *CacheBuilder) *HashListCache {
	hashList := &HashListCache{
		baseCache: baseCache{},
		items:     map[interface{}]*list.Element{},
		lists:     list.New(),
	}
	buildCache(&hashList.baseCache, cb)
	if cb.noLazy {
		c := cron.New()
		_, err := c.AddFunc(fmt.Sprintf("@every %v", hashList.purgeTime), hashList.RemoveOnTime)
		if err != nil {
			log.Error("出现无法恢复的错误，无法生成刷新内容的定时器" + err.Error())
			panic("出现无法恢复的错误，无法生成刷新内容的定时器")
		}
		c.Start()
	}
	return hashList
}

func (h *HashListCache) Set(key, value interface{}) error {
	nowTime := time.Now()
	deleteTime := nowTime.Add(*h.liveTime)
	h.mu.Lock()
	defer h.mu.Unlock()
	element, ok := h.items[key]
	if ok {
		element.Value.(*simpleItem).Value = value
		element.Value.(*simpleItem).DeleteTime = &deleteTime
		return nil
	}
	if h.maxSize <= h.size {
		return errors.New("超出最大容量限制")
	}
	item := &simpleItem{
		Key:        key,
		CreateTime: &nowTime,
		Value:      value,
		DeleteTime: &deleteTime,
		Msg:        nil,
	}
	element = h.lists.PushBack(item)
	h.items[key] = element
	h.size += 1
	return nil
}

func (h *HashListCache) Get(key interface{}) (interface{}, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	element, ok := h.items[key]
	if !ok {
		return nil, errors.New("不存在该值")
	}
	if element.Value.(*simpleItem).DeleteTime.Before(time.Now()) {
		return nil, errors.New("该值已超时")
	}
	deleteTime := time.Now().Add(*h.liveTime)
	element.Value.(*simpleItem).DeleteTime = &deleteTime
	h.lists.Remove(element)
	item := &simpleItem{
		Key:        element.Value.(*simpleItem).Key,
		CreateTime: element.Value.(*simpleItem).CreateTime,
		Value:      element.Value.(*simpleItem).Value,
		DeleteTime: element.Value.(*simpleItem).DeleteTime,
		Msg:        element.Value.(*simpleItem).Msg,
	}
	element = h.lists.PushBack(item)
	h.items[key] = element
	return element.Value.(*simpleItem).Value, nil
}

func (h *HashListCache) GetThroughKeyValue(key interface{}) []KeyAndValue {
	h.mu.RLock()
	defer h.mu.RUnlock()
	valueList := make([]KeyAndValue, 0)
	for k, v := range h.items {
		if reflect.DeepEqual(key, k) {
			if v.Value.(*simpleItem).DeleteTime.Before(time.Now()) {
				continue
			}
			deleteTime := time.Now().Add(*h.liveTime)
			v.Value.(*simpleItem).DeleteTime = &deleteTime
			valueList = append(valueList, KeyAndValue{
				Key:   k,
				Value: v.Value.(*simpleItem).Value,
			})
		}
	}
	return valueList
}

func (h *HashListCache) GetByFunc(key interface{}, fun func(key1, key2 interface{}) bool) []KeyAndValue {
	h.mu.RLock()
	defer h.mu.RUnlock()
	valueList := make([]KeyAndValue, 0)
	for k, v := range h.items {
		if fun(key, k) {
			if v.Value.(*simpleItem).DeleteTime.Before(time.Now()) {
				continue
			}
			deleteTime := time.Now().Add(*h.liveTime)
			v.Value.(*simpleItem).DeleteTime = &deleteTime
			valueList = append(valueList, KeyAndValue{
				Key:   k,
				Value: v.Value.(*simpleItem).Value,
			})
		}
	}
	return valueList
}

func (h *HashListCache) GetAll() []KeyAndValue {
	h.mu.RLock()
	defer h.mu.RUnlock()
	valueList := make([]KeyAndValue, 0)
	for k, v := range h.items {
		if v.Value.(*simpleItem).DeleteTime.Before(time.Now()) {
			continue
		}
		valueList = append(valueList, KeyAndValue{
			Key:   k,
			Value: v.Value.(*simpleItem).Value,
		})
	}
	return valueList
}

func (h *HashListCache) Remove(key interface{}) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.remove(key)
}

func (h *HashListCache) remove(key interface{}) bool {
	element, ok := h.items[key]
	if !ok {
		return false
	}
	if h.deleteFunc != nil {
		err := h.deleteFunc(element.Value.(*simpleItem).Value)
		if err != nil {
			return false
		}
	}
	h.lists.Remove(element)
	delete(h.items, key)
	h.size -= 1
	return true
}

func (h *HashListCache) GetAllKey() []interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()
	valueList := make([]interface{}, 0)
	for k, v := range h.items {
		if v.Value.(*simpleItem).DeleteTime.Before(time.Now()) {
			continue
		}
		valueList = append(valueList, k)
	}
	return valueList
}

func (h *HashListCache) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.size
}

func (h *HashListCache) Has(key interface{}) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.items[key]
	if !ok {
		return false
	}
	return true
}

func (h *HashListCache) HasThroughKeyValue(key interface{}) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for k, v := range h.items {
		if reflect.DeepEqual(key, k) {
			if v.Value.(*simpleItem).DeleteTime.Before(time.Now()) {
				continue
			}
			return true
		}
	}
	return false
}

func (h *HashListCache) HasByFunc(key interface{}, fun func(key1 interface{}, key2 interface{}) bool) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for k, v := range h.items {
		if fun(key, k) {
			if v.Value.(*simpleItem).DeleteTime.Before(time.Now()) {
				continue
			}
			return true
		}
	}
	return false
}

func (h *HashListCache) Range(fun func(k interface{}, v interface{}, result chan interface{}) bool, result chan interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for k, v := range h.items {
		r := fun(k, v.Value.(*simpleItem).Value, result)
		if !r {
			close(result)
			return
		}
	}
	close(result)
}

// 净化cache, 立刻删除过期的值
func (h *HashListCache) Purge() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for k, v := range h.items {
		if v.Value.(*simpleItem).DeleteTime.Before(time.Now()) {
			h.remove(k)
		}
	}
}

func (h *HashListCache) Init() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.size = 0
	h.items = map[interface{}]*list.Element{}
	h.lists = list.New()
}

// RemoveOnTime 删除超时的节点
func (h *HashListCache) RemoveOnTime() {
	fmt.Println("开始删除超时的节点")
	nowTime := time.Now()
	for i := h.lists.Front(); i != nil; i = i.Next() {
		if i.Value.(*simpleItem).DeleteTime.Before(nowTime) {
			h.remove(i.Value.(*simpleItem).Key)
		} else {
			break
		}
	}
}
