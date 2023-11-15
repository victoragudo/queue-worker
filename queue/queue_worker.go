package queue

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

type Worker[T any] struct {
	queue        *Queue[T]
	processor    IProcessor
	wg           sync.WaitGroup
	stopCh       chan struct{}
	processEvery time.Duration
}

func NewQueueWorker[T any](queue *Queue[T], processor IProcessor, processEvery time.Duration) *Worker[T] {
	return &Worker[T]{
		queue:        queue,
		processor:    processor,
		stopCh:       make(chan struct{}),
		processEvery: processEvery,
	}
}

func (w *Worker[T]) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer func() {
			w.closeChannel()
		}()
		var sem = semaphore.NewWeighted(int64(5))
		for {
			select {
			case <-w.stopCh:
				return
			default:
				item, ok := w.queue.TryDequeue()
				if ok {
					fmt.Printf("Dequeue %v\n:", item)
					err := sem.Acquire(context.TODO(), 1)
					if err != nil {
						return
					}
					go func() {
						defer sem.Release(1)
						err := w.processor.Process(item)
						if err != nil {
							println(err)
						}
					}()
				}
				time.Sleep(w.processEvery)
			}
		}
	}()
}

func (w *Worker[T]) Stop() {
	close(w.stopCh)
	w.wg.Wait()
}

func (w *Worker[T]) closeChannel() {
	_, ok := <-w.stopCh
	if ok {
		close(w.stopCh)
	}
}
