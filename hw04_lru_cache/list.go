package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	List // Remove me after realization.
	// Place your code here.
	len int
}

func NewList() List {
	return new(list)
}

type MyList struct {
}

func (ml *list) Len() int {
	return ml.len
}
