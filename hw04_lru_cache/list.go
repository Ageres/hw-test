package hw04lrucache

//----------------------------------------------------------------------------------------------------
// List

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type list struct {
	List  // Remove me after realization.
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v any) *ListItem {
	listItemRef := &ListItem{
		value: v,
	}
	if l.len == 0 {
		l.back = listItemRef
	} else {
		frontListItem := l.front
		frontListItem.prev = listItemRef
		listItemRef.next = frontListItem
	}
	l.front = listItemRef
	l.len++
	return l.back
}

func (l *list) PushBack(v any) *ListItem {
	listItemRef := &ListItem{
		value: v,
	}
	if l.len == 0 {
		l.front = listItemRef
	} else {
		backListItem := l.back
		backListItem.next = listItemRef
		listItemRef.prev = backListItem
	}
	l.back = listItemRef
	l.len++
	return l.back
}

func (l *list) Remove(listItem *ListItem) {
	prev := listItem.prev
	next := listItem.next

	if prev != nil && next != nil {
		prev.next = next
		next.prev = prev
	}

	if prev != nil && next == nil {
		prev.next = nil
		l.back = prev
	}

	if prev == nil && next != nil {
		next.prev = nil
		l.front = next
	}

	if prev == nil && next == nil {
		l.front = nil
		l.back = nil
	}
	listItem.next = nil
	listItem.prev = nil
	l.len--
	/*
		listItemRef := &ListItem{
			Value: v,
		}
		if l.len == 0 {
			l.back = listItemRef
		} else {
			frontListItem := l.front
			frontListItem.prev = listItemRef
			listItemRef.next = frontListItem
		}
		l.front = listItemRef
		l.len++
	*/

}

//----------------------------------------------------------------------------------------------------
// ListItem
type ListItem struct {
	value any
	next  *ListItem
	prev  *ListItem
}

func (li *ListItem) Value() any {
	return li.value
}

func (li *ListItem) Next() *ListItem {
	return li.next
}

func (li *ListItem) Prev() *ListItem {
	return li.prev
}
