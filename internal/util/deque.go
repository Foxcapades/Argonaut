package util

type Deque[T any] interface {
	IsEmpty() bool
	Poll() T
	Offer(value T) bool
	LastIndex() int
}

func NewDeque[T any](initialSize int) Deque[T] {
	return NewCappedDeque[T](-1, initialSize)
}

func NewCappedDeque[T any](maxSize, initialSize int) Deque[T] {
	return &dequeImpl[T]{maxSize: maxSize, initSize: initialSize}
}

type dequeImpl[T any] struct {
	container []T
	realHead  int
	size      int
	maxSize   int
	initSize  int
}

// // // // // // // // // // // // // // // // // // // // // // // // // // //
//
//   INTERNAL API
//
// // // // // // // // // // // // // // // // // // // // // // // // // // //

func (d *dequeImpl[T]) ensureCapacity(minimum int) {
	if !d.canGrowTo(minimum) {
		panic("attempted to grow a sized deque instance to be greater than its max configured size")
	}

	if minimum <= len(d.container) {
		return
	}

	if d.container == nil {
		d.container = make([]T, d.initSize)
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

func (d *dequeImpl[T]) canGrowTo(minimum int) bool {
	return d.maxSize < 0 || minimum <= d.maxSize
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

func (d *dequeImpl[T]) Offer(value T) bool {
	if !d.canGrowTo(d.size + 1) {
		return false
	}

	d.ensureCapacity(d.size + 1)

	d.container[d.internalIndex(d.size)] = value
	d.size++

	return true
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
