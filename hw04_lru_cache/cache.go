package hw04lrucache

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
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value any) bool {
	oldListItem, ok := l.items[key]
	if !ok {
		if l.queue.Len() == l.capacity {
			backListItem := l.queue.Back()
			for k, v := range l.items {
				if v == backListItem {
					delete(l.items, k)
					break
				}
			}
			l.queue.Remove(backListItem)
		}
	} else {
		l.queue.Remove(oldListItem)
	}
	newListItem := l.queue.PushFront(value)
	l.items[key] = newListItem
	return ok
}

func (l *lruCache) Get(key Key) (any, bool) {
	oldListItem, ok := l.items[key]
	if !ok {
		return nil, false
	}
	value := oldListItem.Value()
	l.queue.Remove(oldListItem)
	newListItem := l.queue.PushFront(value)
	l.items[key] = newListItem
	return value, ok
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}
