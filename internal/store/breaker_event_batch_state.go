package store

import "errors"

type BreakerEventBatchStream struct{}

func NewBreakerEventBatchStream() *BreakerEventBatchStream { return &BreakerEventBatchStream{} }

func (s *BreakerEventBatchStream) Start(items []string, failAt int) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error, 1)
	go func() {
		defer close(errs)
		defer close(out)
		for index, item := range items {
			if index == failAt {
				errs <- errors.New("stream interrupted")
				return
			}
			out <- item
		}
	}()
	return out, errs
}
