package lru

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
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
	clearCh  chan string
	openCh   bool
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int, ch chan string) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		clearCh:  ch,
		openCh:   true,
	}
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	listItem, ok := l.items[key]
	if !ok {
		return nil, false
	}

	l.queue.MoveToFront(listItem)

	cacheItem, ok := listItem.Value.(cacheItem)

	if !ok {
		return nil, false
	}

	return cacheItem.value, true
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	listItem, ok := l.items[key]

	if ok {
		listItem.Value = createCacheItem(key, value)

		l.queue.MoveToFront(listItem)

		return true
	}

	if l.queue.Len() == l.capacity {
		removeListItem := l.queue.Back()
		removeCacheItem, ok := removeListItem.Value.(cacheItem)
		if ok {
			delete(l.items, removeCacheItem.key)
		}

		l.clearCh <- string(removeCacheItem.key)
		l.queue.Remove(removeListItem)
	}

	newListItem := l.queue.PushFront(createCacheItem(key, value))

	l.items[key] = newListItem

	return false
}

func createCacheItem(key Key, value interface{}) cacheItem {
	return cacheItem{
		key:   key,
		value: value,
	}
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}
