package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       *sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mu:       new(sync.Mutex),
	}
}

func (l *lruCache) Set(key Key, value any) bool {
	_, ok := l.items[key]
	if !ok {
		if l.queue.Len() == l.capacity {
			backListItem := l.queue.Back()
			for k, v := range l.items {
				if v == backListItem.Value() {
					delete(l.items, k) // а если добавляли несколько раз один и тот же элемент с разными ключами?
					break
				}
			}
			l.queue.Remove(backListItem)
		}
	}
	newListItem := l.queue.PushFront(value)
	l.items[key] = newListItem
	return ok
}

func (l *lruCache) Get(key Key) (any, bool) {
	oldListItem, ok := l.items[key]
	if !ok {
		return nil, false
	} else {
		value := oldListItem.Value()
		l.queue.Remove(oldListItem)
		newListItem := l.queue.PushFront(value)
		l.items[key] = newListItem
		return value, ok
	}
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}
