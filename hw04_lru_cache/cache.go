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
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.Mutex
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	if lc.capacity == 0 {
		return false
	}
	// Ключ существует, обновляем значение
	if item, ok := lc.items[key]; ok {
		item.Value = value
		item.Key = key
		lc.queue.MoveToFront(item)
		log.Printf("Key %s updated to %v\n", key, value)
		return true
	}
	// Добавление нового элемента
	lc.queue.PushFront(value, key)
	lc.items[key] = lc.queue.Front()

	// Удаление превышения кэша
	if lc.queue.Len() > lc.capacity {
		oldest := lc.queue.Back()
		oldKey := oldest.Key
		lc.queue.Remove(oldest)
		delete(lc.items, oldKey)
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
	item.Key = key
	lc.queue.MoveToFront(item)
	return item.Value, true
}

func (lc *lruCache) Clear() {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	lc.items = map[Key]*ListItem{}
	lc.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mutex:    sync.Mutex{},
	}
}
