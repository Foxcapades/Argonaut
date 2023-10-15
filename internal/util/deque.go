package util

import (
	"fmt"

	"github.com/Foxcapades/Argonaut/v1/internal/xmath"
)

type Deque[T any] interface {
	IsEmpty() bool
	Poll() T
	Offer(value T) bool
	AddAll(values []T) int
}

// NewDeque creates a new, empty DoubleEndedQueue instance with the given
// initial capacity.
func NewDeque[T any](initialSize int) Deque[T] {
	return NewCappedDeque[T](-1, initialSize)
}

// NewCappedDeque creates a new DoubleEndedQueue with the given initial size and
// max capacity values.
func NewCappedDeque[T any](maxSize, initialSize int) Deque[T] {
	return &deque[T]{maxSize: maxSize, initSize: initialSize}
}

type deque[T any] struct {
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

func (d *deque[T]) compact() {
	d.copyElements(len(d.container))
}

func (d *deque[T]) trimToSize() {
	d.copyElements(d.size)
}

func (d *deque[T]) isInline() bool {
	return d.realHead <= d.LastIndex()
}

func (d *deque[T]) ensureCapacity(minimum int) {
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

func (d *deque[T]) vei(i int) int {
	if i < 0 || i > d.size-1 {
		panic("index out of bounds")
	} else {
		return d.internalIndex(i)
	}
}

func (d *deque[T]) positiveMod(i int) int {
	if i >= len(d.container) {
		return i - len(d.container)
	} else {
		return i
	}
}

func (d *deque[T]) negativeMod(i int) int {
	if i < 0 {
		return i + len(d.container)
	} else {
		return i
	}
}

func (d *deque[T]) internalIndex(i int) int {
	return d.positiveMod(d.realHead + i)
}

func (d *deque[T]) incremented(i int) int {
	if i == len(d.container)-1 {
		return 0
	} else {
		return i + 1
	}
}

func (d *deque[T]) decremented(i int) int {
	if i < 1 {
		return len(d.container) - 1
	} else {
		return i - 1
	}
}

func (d *deque[T]) copyElements(capacity int) {
	newContainer := make([]T, capacity)
	copy(newContainer, d.container[d.realHead:])
	copy(newContainer[len(d.container)-d.realHead:], d.container[:d.realHead])
	d.realHead = 0
	d.container = newContainer
}

func (d *deque[T]) newCap(old, min int) int {
	return xmath.Max(old+(old>>1), min)
}

func (d *deque[T]) canGrowTo(minimum int) bool {
	return d.maxSize < 0 || minimum <= d.maxSize
}

// // // // // // // // // // // // // // // // // // // // // // // // // // //
//
//   PUBLIC API
//
// // // // // // // // // // // // // // // // // // // // // // // // // // //

func (d *deque[T]) AddAll(values []T) int {
	if d.canGrowTo(d.size + len(values)) {
		d.ensureCapacity(d.size + len(values))
	} else {
		d.ensureCapacity(d.maxSize)
	}

	count := 0

	for i := range values {
		if d.Offer(values[i]) {
			count++
		} else {
			break
		}
	}

	return count
}

func (d *deque[T]) IsEmpty() bool {
	return d.size == 0
}

func (d *deque[T]) LastIndex() int {
	return d.size - 1
}

func (d *deque[T]) Offer(value T) bool {
	if !d.canGrowTo(d.size + 1) {
		return false
	}

	d.ensureCapacity(d.size + 1)

	d.container[d.internalIndex(d.size)] = value
	d.size++

	return true
}

func (d *deque[T]) Poll() T {
	if d.size == 0 {
		panic("no such element")
	}

	out := d.container[d.realHead]

	d.realHead = d.incremented(d.realHead)
	d.size--

	return out
}

func (d *deque[T]) Size() int {
	return d.size
}

func (d *deque[T]) String() string {
	if d.maxSize > -1 {
		return fmt.Sprintf("Deque{Size: %d, Capacity: %d, MaxSize: %d}", d.size, len(d.container), d.maxSize)
	} else {
		return fmt.Sprintf("Deque{Size: %d, Capacity: %d}", d.size, len(d.container))
	}
}
