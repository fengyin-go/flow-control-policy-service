package store

import "errors"

type DecisionBatchStream struct{}

func NewDecisionBatchStream() *DecisionBatchStream { return &DecisionBatchStream{} }

func (s *DecisionBatchStream) Start(items []string, failAt int) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error, 1)
	go func() {

		defer close(errs)
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
