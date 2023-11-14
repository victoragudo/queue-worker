package queue

import (
	"sync"
)

type Queue[T any] struct {
	items chan T
	size  int
	mu    sync.RWMutex
}

func NewQueue[T any](size int) *Queue[T] {
	return &Queue[T]{
		items: make(chan T, size),
		size:  0,
	}
}

func (q *Queue[T]) Enqueue(item T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	select {
	case q.items <- item:
		q.size++
		return true
	default:
		return false
	}
}

func (q *Queue[T]) TryDequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	select {
	case item := <-q.items:
		q.size--
		return item, true
	default:
		zeroValue := new(T)
		return *zeroValue, false
	}
}

func (q *Queue[T]) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.size
}

func (q *Queue[T]) Capacity() int {
	return cap(q.items)
}
