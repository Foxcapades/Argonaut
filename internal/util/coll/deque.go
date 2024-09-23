package coll

type Deque[T any] interface {
	IsEmpty() bool
	Poll() T
	PushBack(value T)
	PushFront(value T)
	LastIndex() int
}

func NewDeque[T any](initialSize int) Deque[T] {
	return &dequeImpl[T]{container: make([]T, initialSize)}
}

type dequeImpl[T any] struct {
	container []T
	realHead  int
	size      int
}

// // // // // // // // // // // // // // // // // // // // // // // // // // //
//
//   INTERNAL API
//
// // // // // // // // // // // // // // // // // // // // // // // // // // //

func (d *dequeImpl[T]) ensureCapacity(minimum int) {
	if minimum <= len(d.container) {
		return
	}

	d.copyElements(d.newCap(len(d.container), minimum))
}

func (d *dequeImpl[T]) positiveMod(i int) int {
	if i >= len(d.container) {
		return i - len(d.container)
	}

	return i
}

func (d *dequeImpl[T]) internalIndex(i int) int {
	return d.positiveMod(d.realHead + i)
}

func (d *dequeImpl[T]) incremented(i int) int {
	if i == len(d.container)-1 {
		return 0
	}

	return i + 1
}

func (d *dequeImpl[T]) decremented(i int) int {
	if i == 0 {
		return len(d.container) - 1
	}

	return i - 1
}

func (d *dequeImpl[T]) copyElements(capacity int) {
	newContainer := make([]T, capacity)
	copy(newContainer, d.container[d.realHead:])
	copy(newContainer[len(d.container)-d.realHead:], d.container[:d.realHead])
	d.realHead = 0
	d.container = newContainer
}

func (d *dequeImpl[T]) newCap(old, min int) int {
	return max(old+(old>>1), min)
}

// // // // // // // // // // // // // // // // // // // // // // // // // // //
//
//   PUBLIC API
//
// // // // // // // // // // // // // // // // // // // // // // // // // // //

func (d *dequeImpl[T]) IsEmpty() bool {
	return d.size == 0
}

func (d *dequeImpl[T]) LastIndex() int {
	return d.size - 1
}

func (d *dequeImpl[T]) PushFront(value T) {
	d.ensureCapacity(d.size + 1)

	d.realHead = d.decremented(d.realHead)
	d.container[d.realHead] = value
	d.size++
}

func (d *dequeImpl[T]) PushBack(value T) {
	d.ensureCapacity(d.size + 1)

	d.container[d.internalIndex(d.size)] = value
	d.size++
}

func (d *dequeImpl[T]) Poll() T {
	if d.size == 0 {
		panic("no such element")
	}
	var t T

	out := d.container[d.realHead]
	d.container[d.realHead] = t

	d.realHead = d.incremented(d.realHead)
	d.size--

	return out
}
