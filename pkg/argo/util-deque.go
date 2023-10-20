package argo

type deque[T any] interface {
	IsEmpty() bool
	Poll() T
	Offer(value T) bool
}

func newDeque[T any](initialSize int) deque[T] {
	return newCappedDeque[T](-1, initialSize)
}

func newCappedDeque[T any](maxSize, initialSize int) deque[T] {
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

func (d *dequeImpl[T]) compact() {
	d.copyElements(len(d.container))
}

func (d *dequeImpl[T]) trimToSize() {
	d.copyElements(d.size)
}

func (d *dequeImpl[T]) isInline() bool {
	return d.realHead <= d.LastIndex()
}

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

func (d *dequeImpl[T]) vei(i int) int {
	if i < 0 || i > d.size-1 {
		panic("index out of bounds")
	} else {
		return d.internalIndex(i)
	}
}

func (d *dequeImpl[T]) positiveMod(i int) int {
	if i >= len(d.container) {
		return i - len(d.container)
	} else {
		return i
	}
}

func (d *dequeImpl[T]) negativeMod(i int) int {
	if i < 0 {
		return i + len(d.container)
	} else {
		return i
	}
}

func (d *dequeImpl[T]) internalIndex(i int) int {
	return d.positiveMod(d.realHead + i)
}

func (d *dequeImpl[T]) incremented(i int) int {
	if i == len(d.container)-1 {
		return 0
	} else {
		return i + 1
	}
}

func (d *dequeImpl[T]) decremented(i int) int {
	if i < 1 {
		return len(d.container) - 1
	} else {
		return i - 1
	}
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
