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
	var ok bool

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer l.mu.Unlock()
		l.mu.Lock()

		var oldListItem *ListItem
		oldListItem, ok = l.items[key]
		if !ok {
			if l.queue.Len() == l.capacity {
				backListItem := l.queue.Back()
				delete(l.items, getKey(backListItem))
				l.queue.Remove(backListItem)
			}
		} else {
			l.queue.Remove(oldListItem)
		}

		newCacheItem := &cacheItem{
			Key:   key,
			Value: value,
		}
		newListItem := l.queue.PushFront(newCacheItem)
		l.items[key] = newListItem
	}()

	wg.Wait()

	return ok
}

func (l *lruCache) Get(key Key) (any, bool) {
	var value any
	var ok bool

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer l.mu.Unlock()
		l.mu.Lock()

		var oldListItem *ListItem
		oldListItem, ok = l.items[key]
		if !ok {
			return
		}

		value = getValue(oldListItem)
		l.queue.Remove(oldListItem)

		newCacheItem := &cacheItem{
			Key:   key,
			Value: value,
		}

		newListItem := l.queue.PushFront(newCacheItem)
		l.items[key] = newListItem
	}()

	wg.Wait()

	return value, ok
}

func (l *lruCache) Clear() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer l.mu.Unlock()
		l.mu.Lock()

		l.queue = NewList()
		l.items = make(map[Key]*ListItem, l.capacity)
	}()

	wg.Wait()
}

//----------------------------------------------------------------------------------------------------
// вспомогательные структуры и методы

// элемент, хранящийся в очереди и словаре в составе ListItem.
type cacheItem struct {
	Key   Key
	Value any
}

func getKey(i *ListItem) Key {
	return i.Value().(*cacheItem).Key
}

func getValue(i *ListItem) any {
	return i.Value().(*cacheItem).Value
}
