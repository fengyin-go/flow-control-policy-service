package service

import "flowcontrol/internal/store"

type QuotaConsumeBatchFlow struct {
	stream *store.QuotaConsumeBatchStream
}

func NewQuotaConsumeBatchFlow(stream *store.QuotaConsumeBatchStream) *QuotaConsumeBatchFlow {
	return &QuotaConsumeBatchFlow{stream: stream}
}

func (f *QuotaConsumeBatchFlow) Collect(items []string, failAt int) ([]string, error) {
	out, errs := f.stream.Start(items, failAt)
	values := make([]string, 0)
	for value := range out {
		values = append(values, value)
	}
	return values, <-errs
}
