package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{}) *listItem
	PushBack(v interface{}) *listItem
	Remove(i *listItem)
	MoveToFront(i *listItem)
}

type listItem struct {
	Value interface{}
	Next  *listItem
	Prev  *listItem
}

type list struct {
	front *listItem
	back  *listItem
	len   int
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *listItem {
	return l.front
}

func (l *list) Back() *listItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *listItem {
	newItem := listItem{v, l.front, nil}
	if l.Len() != 0 {
		l.front.Prev = &newItem
	}
	l.front = &newItem
	l.len++
	if l.Len() == 1 {
		l.back = l.front
	}
	return l.front
}

func (l *list) PushBack(v interface{}) *listItem {
	newItem := listItem{v, nil, l.back}
	if l.Len() != 0 {
		l.back.Next = &newItem
	}
	l.back = &newItem
	l.len++
	if l.Len() == 1 {
		l.front = l.back
	}
	return l.back
}

func (l *list) Remove(i *listItem) {
	switch {
	case l.Front() == l.Back():
		l.front = nil
		l.back = nil
	case l.Front() == i:
		l.front = l.front.Next
		l.front.Prev = nil
	case l.Back() == i:
		l.back = l.back.Prev
		l.back.Next = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *listItem) {
	l.Remove(i)
	if l.Len() != 0 {
		l.front.Prev = i
		i.Next = l.front
	}
	l.front = i
	i.Prev = nil
	l.len++
	if l.Len() == 1 {
		l.back = l.front
	}
}
