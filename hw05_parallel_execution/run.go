package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(taskCh <-chan Task, doneCh <-chan struct{}, errCh chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-doneCh:
			return
		default:
		}

		select {
		case t, ok := <-taskCh:
			if !ok {
				return
			} else if err := t(); err != nil {
				errCh <- struct{}{}
			}
		case <-doneCh:
			return
		}
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	nErr := 0
	taskChLimit := 5
	taskCh := make(chan Task, taskChLimit)
	doneCh := make(chan struct{})
	errCh := make(chan struct{})
	wgWorker := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wgWorker.Add(1)
		go worker(taskCh, doneCh, errCh, &wgWorker)
	}

	wgCounter := sync.WaitGroup{}
	wgCounter.Add(1)
	go func() {
		defer wgCounter.Done()
		for {
			if _, ok := <-errCh; !ok {
				return
			}
			nErr++
			if nErr == m {
				close(doneCh)
			}
		}
	}()

	for i := 0; i < len(tasks); i++ {
		select {
		case taskCh <- tasks[i]:
		case <-doneCh:
			break
		}
	}
	close(taskCh)

	wgWorker.Wait()
	close(errCh)
	wgCounter.Wait()

	if nErr >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
