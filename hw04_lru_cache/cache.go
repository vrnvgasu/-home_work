package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity     int
	queue        List
	items        map[Key]*ListItem
	itemsReverse map[*ListItem]Key
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	item, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(item)
		front := l.queue.Front()
		front.Value = value
		l.items[key] = item
		l.itemsReverse[item] = key

		return true
	}

	if l.capacity == l.queue.Len() {
		back := l.queue.Back()

		oldKey := l.itemsReverse[back]
		delete(l.items, oldKey)
		delete(l.itemsReverse, back)

		l.queue.Remove(back)
	}

	newItem := l.queue.PushFront(value)
	l.items[key] = newItem
	l.itemsReverse[newItem] = key

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := l.items[key]
	if !ok {
		return nil, false
	}

	l.queue.MoveToFront(item)

	return item.Value, true
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
	l.itemsReverse = make(map[*ListItem]Key, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity:     capacity,
		queue:        NewList(),
		items:        make(map[Key]*ListItem, capacity),
		itemsReverse: make(map[*ListItem]Key, capacity),
	}
}
