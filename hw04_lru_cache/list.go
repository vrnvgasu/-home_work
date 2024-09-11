package hw04lrucache

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
	frontItem   *ListItem
	backItem    *ListItem
	listItemMap map[*ListItem]struct{}
}

func (l *list) Len() int {
	return len(l.listItemMap)
}

func (l *list) Front() *ListItem {
	return l.frontItem
}

func (l *list) Back() *ListItem {
	return l.backItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.Len() == 0 {
		l.frontItem = &ListItem{v, nil, nil}
		l.backItem = l.frontItem
		l.listItemMap[l.frontItem] = struct{}{}

		return l.frontItem
	}

	oldFrontItem := l.frontItem
	l.frontItem = &ListItem{
		Value: v,
		Next:  oldFrontItem,
		Prev:  nil,
	}
	oldFrontItem.Prev = l.frontItem
	l.listItemMap[l.frontItem] = struct{}{}

	return l.frontItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.Len() == 0 {
		l.backItem = &ListItem{v, nil, nil}
		l.frontItem = l.backItem
		l.listItemMap[l.backItem] = struct{}{}

		return l.backItem
	}

	oldBackItem := l.backItem
	l.backItem = &ListItem{
		Value: v,
		Next:  nil,
		Prev:  oldBackItem,
	}
	oldBackItem.Next = l.backItem
	l.listItemMap[l.backItem] = struct{}{}

	return l.backItem
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	prev := i.Prev
	next := i.Next

	if prev != nil {
		prev.Next = next
	} else {
		l.frontItem = next
	}

	if next != nil {
		next.Prev = prev
	} else {
		l.backItem = prev
	}

	delete(l.listItemMap, i)

	if len(l.listItemMap) == 0 {
		l.frontItem = nil
		l.backItem = nil
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}
	if l.frontItem == i {
		return
	}

	prev := i.Prev
	next := i.Next

	// уверены, что это не первый элемент
	prev.Next = next

	if next != nil {
		next.Prev = prev
	} else {
		l.backItem = prev
	}

	i.Next = l.frontItem
	l.frontItem.Prev = i
	l.frontItem = i
}

func NewList() List {
	return &list{listItemMap: make(map[*ListItem]struct{})}
}
