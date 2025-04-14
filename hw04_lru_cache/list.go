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
	if l.len == 0 {

	} else {

	}

	return l.back
}

//----------------------------------------------------------------------------------------------------
// ListItem
type ListItem struct {
	Value any
	next  *ListItem
	prev  *ListItem
}

func (li *ListItem) Next() *ListItem {
	return li.next
}

func (li *ListItem) Prev() *ListItem {
	return li.prev
}
