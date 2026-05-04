package hw04lrucache

import (
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
	mutex    sync.RWMutex
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
		lc.queue.MoveToFront(item)
		return true // Элемент уже был в кэше
	}
	// Добавление нового элемента
	lc.queue.PushFront(value)
	lc.items[key] = lc.queue.Front()

	// Удаление превышения кэша
	if lc.queue.Len() > lc.capacity {
		t := lc.queue.Back()
		lc.queue.Remove(t)
		for k, v := range lc.items {
			if v.Value == t.Value {
				delete(lc.items, k)
			}
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
	head := lc.queue.Front()
	tail := lc.queue.Back()
	head.Next = nil
	tail.Prev = nil
	lc.queue.Remove(head)
	lc.queue.Remove(tail)
	lc.queue.Len()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mutex:    sync.RWMutex{},
	}
}
