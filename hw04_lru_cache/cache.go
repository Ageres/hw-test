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
			delete(l.items, getKey(backListItem))
			l.queue.Remove(backListItem)
		}
	} else {
		l.queue.Remove(oldListItem)
	}

	newCasheItem := &CacheItem{
		Key:   key,
		Value: value,
	}
	newListItem := l.queue.PushFront(newCasheItem)
	l.items[key] = newListItem
	return ok
}

func (l *lruCache) Get(key Key) (any, bool) {
	oldListItem, ok := l.items[key]
	if !ok {
		return nil, false
	}

	value := getValue(oldListItem)
	l.queue.Remove(oldListItem)

	newCasheItem := &CacheItem{
		Key:   key,
		Value: value,
	}

	newListItem := l.queue.PushFront(newCasheItem)
	l.items[key] = newListItem
	return value, ok
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

//----------------------------------------------------------------------------------------------------
// CasheItem

// элемент, хранящийся в очереди и словаре в составе ListItem.
type CacheItem struct {
	Key   Key
	Value any
}

func getKey(i *ListItem) Key {
	return i.Value().(*CacheItem).Key
}

func getValue(i *ListItem) any {
	return i.Value().(*CacheItem).Value
}
