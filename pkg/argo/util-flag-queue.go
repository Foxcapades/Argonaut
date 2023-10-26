package argo

import "github.com/Foxcapades/Argonaut/internal/util"

func newFlagQueue() flagQueue {
	return flagQueue{
		ordered:  make([]util.Pair[byte, string], 0, 4),
		distinct: make(map[util.Pair[byte, string]]Flag, 4),
	}
}

type flagQueue struct {
	ordered  []util.Pair[byte, string]
	distinct map[util.Pair[byte, string]]Flag
}

func (q *flagQueue) append(f Flag) {
	var key = util.Pair[byte, string]{L: f.ShortForm(), R: f.LongForm()}
	if _, ok := q.distinct[key]; !ok {
		q.distinct[key] = f
		q.ordered = append(q.ordered, key)
	}
}

func (q *flagQueue) iterator() flagQueueIterator {
	return flagQueueIterator{queue: q}
}

type flagQueueIterator struct {
	queue    *flagQueue
	position int
}

func (i *flagQueueIterator) hasNext() bool {
	return i.position < len(i.queue.ordered)
}

func (i *flagQueueIterator) next() Flag {
	if i.position >= len(i.queue.ordered) {
		panic("no such element")
	}

	out := i.queue.distinct[i.queue.ordered[i.position]]
	i.position++

	return out
}
