package service

import "flowcontrol/internal/store"

type BreakerEventBatchFlow struct {
	stream *store.BreakerEventBatchStream
}

func NewBreakerEventBatchFlow(stream *store.BreakerEventBatchStream) *BreakerEventBatchFlow {
	return &BreakerEventBatchFlow{stream: stream}
}

func (f *BreakerEventBatchFlow) Collect(items []string, failAt int) ([]string, error) {
	out, errs := f.stream.Start(items, failAt)
	values := make([]string, 0)
	for value := range out {
		values = append(values, value)
	}
	return values, <-errs
}
