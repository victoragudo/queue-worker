package queue

import (
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
		for {
			select {
			case <-w.stopCh:
				return
			default:
				item, ok := w.queue.TryDequeue()
				if ok {
					err := w.processor.Process(item)
					if err != nil {
						println(err)
					}
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
