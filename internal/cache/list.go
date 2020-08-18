package cache //nolint:golint,stylecheck

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
	if l == nil {
		return 0
	}

	return l.len
}

func (l *list) Front() *listItem {
	return l.front
}

func (l *list) Back() *listItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *listItem { //nolint:dupl
	itm := &listItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
	if l.len == 0 { // empty list
		l.front = itm
		l.len++

		return l.front
	}
	if l.front == nil { // one element in "back"
		itm.Prev = l.back
		l.front = itm
		l.back.Next = l.front
		l.len++

		return l.front
	}
	if l.back == nil { // one element in "front"
		l.back = l.front
		itm.Prev = l.back
		l.front = itm
		l.back.Next = l.front
		l.len++

		return l.front
	}
	itm.Prev = l.front
	l.front.Next = itm
	l.front = itm
	l.len++

	return l.front
}

func (l *list) PushBack(v interface{}) *listItem { //nolint:dupl
	itm := &listItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
	if l.len == 0 { // empty list
		l.back = itm
		l.len++

		return l.back
	}
	if l.back == nil { // one element in "front"
		itm.Next = l.front
		l.back = itm
		l.front.Prev = l.back
		l.len++

		return l.back
	}

	if l.front == nil { // one element in "back"
		l.front = l.back
		itm.Next = l.front
		l.back = itm
		l.front.Prev = l.back
		l.len++

		return l.back
	}
	itm.Next = l.back
	l.back.Prev = itm
	l.back = itm
	l.len++

	return l.back
}

func (l *list) Remove(i *listItem) {
	if l.len == 0 { // empty list
		return
	}
	if i == l.front {
		if l.front.Prev != nil {
			l.front = l.front.Prev
			if l.front == l.back {
				l.back = nil
			}
			l.front.Next = nil
		} else {
			l.front = nil
		}
		l.len--

		return
	}
	if i == l.back {
		if l.back.Next != nil {
			l.back = l.back.Next
			if l.back == l.front {
				l.front = nil
			}
			l.back.Prev = nil
		} else {
			l.back = nil
		}
		l.len--

		return
	}
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	l.len--
}

func (l *list) MoveToFront(i *listItem) {
	if l.len == 0 { // empty list
		return
	}
	if i == l.front {
		return
	}
	if l.front == nil { // one element in "back"
		l.front = i
		l.back = nil

		return
	}
	if i == l.back {
		l.back = l.back.Next
		l.back.Prev = nil
		l.front.Next = i
		i.Prev = l.front
		i.Next = nil
		l.front = i

		return
	}
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	l.front.Next = i
	i.Prev = l.front
	i.Next = nil
	l.front = i
}
