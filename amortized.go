package list

// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type AmortizedList[T any] struct {
	root   Element[T] // sentinel list element, only &root, root.prev, and root.next are used
	len    int        // current list length excluding (this) sentinel element
	memory *List[T]
}

// Init initializes or clears list l.
func (l *AmortizedList[T]) Init() *AmortizedList[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	l.memory = New[T]()

	return l
}

// New returns an initialized list.
func NewAmortized[T any]() *AmortizedList[T] { return new(AmortizedList[T]).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *AmortizedList[T]) Len() int { return l.len }

// Front returns the first element of list l or nil if the list is empty.
func (l *AmortizedList[T]) Front() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *AmortizedList[T]) Back() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero List value.
func (l *AmortizedList[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *AmortizedList[T]) insert(e, at *Element[T]) *Element[T] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.alist = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *AmortizedList[T]) insertValue(v T, at *Element[T]) *Element[T] {
	lastSlot := l.memory.Back()

	if lastSlot == nil {
		return l.insert(&Element[T]{Value: v}, at)
	}

	l.memory.Remove(lastSlot)
	lastSlot.Value = v

	return l.insert(lastSlot, at)
}

// remove removes e from its list, decrements l.len
func (l *AmortizedList[T]) remove(e *Element[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--

	l.memory.insert(e, l.memory.Back())
}

// move moves e to next to at.
func (l *AmortizedList[T]) move(e, at *Element[T]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *AmortizedList[T]) Remove(e *Element[T]) T {
	if e.alist == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *AmortizedList[T]) PushFront(v T) *Element[T] {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *AmortizedList[T]) PushBack(v T) *Element[T] {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *AmortizedList[T]) InsertBefore(v T, mark *Element[T]) *Element[T] {
	if mark.alist != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *AmortizedList[T]) InsertAfter(v T, mark *Element[T]) *Element[T] {
	if mark.alist != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *AmortizedList[T]) MoveToFront(e *Element[T]) {
	if e.alist != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *AmortizedList[T]) MoveToBack(e *Element[T]) {
	if e.alist != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *AmortizedList[T]) MoveBefore(e, mark *Element[T]) {
	if e.alist != l || e == mark || mark.alist != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *AmortizedList[T]) MoveAfter(e, mark *Element[T]) {
	if e.alist != l || e == mark || mark.alist != l {
		return
	}
	l.move(e, mark)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *AmortizedList[T]) PushBackList(other *AmortizedList[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *AmortizedList[T]) PushFrontList(other *AmortizedList[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}
