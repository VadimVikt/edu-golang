package hw04lrucache

import (
	"log"
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity   int
	queue      List
	items      map[Key]*ListItem
	valueToKey map[interface{}]Key // инвертированная карта: значение -> ключ
	mutex      sync.Mutex
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	if lc.capacity == 0 {
		return false
	}
	// Ключ существует, обновляем значение
	if item, ok := lc.items[key]; ok {
		// Удаляем старое значение из valueToKey
		delete(lc.valueToKey, item.Value)
		item.Value = value
		lc.queue.MoveToFront(item)
		// Добавляем новое значение в valueToKey
		lc.valueToKey[value] = key
		log.Printf("Key %s updated to %v\n", key, value)
		return true
	}
	// Добавление нового элемента
	lc.queue.PushFront(value)
	lc.items[key] = lc.queue.Front()
	lc.valueToKey[value] = key

	// Удаление превышения кэша
	if lc.queue.Len() > lc.capacity {
		oldest := lc.queue.Back()
		lc.queue.Remove(oldest)
		// Находим ключ по значению через valueToKey
		if k, ok := lc.valueToKey[oldest.Value]; ok {
			delete(lc.items, k)
			delete(lc.valueToKey, oldest.Value)
		}
	}
	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	if lc.items[key] == nil {
		return nil, false
	}

	item := lc.items[key]
	lc.queue.MoveToFront(item)
	return item.Value, true
}

func (lc *lruCache) Clear() {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	lc.items = map[Key]*ListItem{}
	lc.valueToKey = map[interface{}]Key{}
	lc.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity:   capacity,
		queue:      NewList(),
		items:      make(map[Key]*ListItem, capacity),
		valueToKey: make(map[interface{}]Key, capacity),
		mutex:      sync.Mutex{},
	}
}
