package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() List {
	return &list{}
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

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := ListItem{v, l.front, nil}
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

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := ListItem{v, nil, l.back}
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

func (l *list) Remove(i *ListItem) {
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

func (l *list) MoveToFront(i *ListItem) {
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
