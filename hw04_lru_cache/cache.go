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
	keys     map[*ListItem]Key
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keys:     make(map[*ListItem]Key, capacity),
	}
}

func (l *lruCache) Set(key Key, value any) bool {
	oldListItem, ok := l.items[key]
	if !ok {
		if l.queue.Len() == l.capacity {
			backListItem := l.queue.Back()
			backKey := l.keys[backListItem]
			delete(l.items, backKey)
			delete(l.keys, backListItem)
			l.queue.Remove(backListItem)
		}
	} else {
		delete(l.keys, oldListItem)
		l.queue.Remove(oldListItem)
	}

	newListItem := l.queue.PushFront(value)
	l.items[key] = newListItem
	l.keys[newListItem] = key
	return ok
}

func (l *lruCache) Get(key Key) (any, bool) {
	oldListItem, ok := l.items[key]
	if !ok {
		return nil, false
	}

	delete(l.keys, oldListItem)
	value := oldListItem.Value()
	l.queue.Remove(oldListItem)

	newListItem := l.queue.PushFront(value)
	l.items[key] = newListItem
	l.keys[newListItem] = key
	return value, ok
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
	l.keys = make(map[*ListItem]Key, l.capacity)
}
