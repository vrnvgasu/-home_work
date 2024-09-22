package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var wg sync.WaitGroup
	ch := make(chan Task)
	defer func() {
		close(ch)
		wg.Wait()
	}()

	var counter int32

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				task, ok := <-ch
				if !ok {
					return
				}
				if err := task(); err != nil {
					atomic.AddInt32(&counter, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&counter) >= int32(m) {
			return ErrErrorsLimitExceeded
		}

		ch <- task
	}

	return nil
}
