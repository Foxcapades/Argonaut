package util_test

import (
	"testing"

	"github.com/Foxcapades/Argonaut/internal/util"
)

func TestDequeImpl_IsEmpty(t *testing.T) {
	if !util.NewDeque[int](10).IsEmpty() {
		t.Error("expected new dequeue to be empty but it wasn't")
	}

	deq := util.NewDeque[int](10)
	for i := 0; i < 12; i++ {
		deq.Offer(i)
	}

	for i := 0; i < 12; i++ {
		deq.Poll()
	}

	if !deq.IsEmpty() {
		t.Error("expected deque to be empty but it wasn't")
	}
}

func TestDequeImpl_LastIndex(t *testing.T) {
	deq := util.NewDeque[int](10)

	if deq.LastIndex() != -1 {
		t.Errorf("expected lastIndex to be -1 but it was %d", deq.LastIndex())
	}

	for i := 0; i < 10; i++ {
		deq.Offer(i)
	}

	if deq.LastIndex() != 9 {
		t.Error("expected lastIndex to be 9 but it wasn't")
	}
}

func TestDequeImpl_Offer(t *testing.T) {
	deq := util.NewDeque[int](10)

	for i := 0; i < 10; i++ {
		deq.Offer(i)
	}

	for i := 0; i < 5; i++ {
		deq.Poll()
	}

	for i := 0; i < 10; i++ {
		deq.Offer(i)
	}
}

func TestDequeImpl_Poll(t *testing.T) {
	defer func() { recover() }()

	deq := util.NewDeque[int](1)

	deq.Poll()

	t.Error("expected Deque.Poll to panic but it didn't")
}
