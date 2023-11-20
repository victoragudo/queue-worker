package main

import (
	"errors"
	"fmt"
	"queue-worker/queue"
	"time"
)

type Email struct {
	Message string
	t       time.Time
}

func main() {
	emailQueue := queue.NewQueue[Email](1 << 10)
	emailProcessor := queue.NewProcessor(func(item any) error {
		v, ok := item.(Email)
		if ok {
			fmt.Printf("%v: %v\n", v.t.String(), v.Message)
			time.Sleep(time.Second * 10)
			return nil
		}
		return errors.New("unable to cast type")
	})

	worker := queue.NewQueueWorker(emailQueue, emailProcessor, time.Millisecond*300)

	worker.Start()

	go func() {
		for {
			emailQueue.Enqueue(Email{Message: "Otro mensaje...", t: time.Now()})
			time.Sleep(time.Millisecond * 100)
		}
	}()

	select {}
	//worker.Stop()
}
