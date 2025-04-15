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
	ch := make(chan bool)

	go func(key Key, value any) {
		defer close(ch)
		defer l.mu.Unlock()
		l.mu.Lock()

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

		newCacheItem := &cacheItem{
			Key:   key,
			Value: value,
		}
		newListItem := l.queue.PushFront(newCacheItem)
		l.items[key] = newListItem
		ch <- ok
	}(key, value)

	return <-ch
}

func (l *lruCache) Get(key Key) (any, bool) {
	ch := make(chan goGetResp)

	go func(key Key) {
		defer close(ch)
		defer l.mu.Unlock()
		l.mu.Lock()

		oldListItem, ok := l.items[key]
		if !ok {
			ch <- goGetResp{
				Ok: false,
			}
			return
		}

		value := getValue(oldListItem)
		l.queue.Remove(oldListItem)

		newCacheItem := &cacheItem{
			Key:   key,
			Value: value,
		}

		newListItem := l.queue.PushFront(newCacheItem)
		l.items[key] = newListItem
		ch <- goGetResp{
			Value: value,
			Ok:    ok,
		}
	}(key)

	resp := <-ch
	return resp.Value, resp.Ok
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

// структура для возврата ответа из горутины в методе Get.
type goGetResp struct {
	Value any
	Ok    bool
}
