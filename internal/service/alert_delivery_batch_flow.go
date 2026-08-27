package service

import "flowcontrol/internal/store"

type AlertDeliveryBatchFlow struct {
	stream *store.AlertDeliveryBatchStream
}

func NewAlertDeliveryBatchFlow(stream *store.AlertDeliveryBatchStream) *AlertDeliveryBatchFlow {
	return &AlertDeliveryBatchFlow{stream: stream}
}

func (f *AlertDeliveryBatchFlow) Collect(items []string, failAt int) ([]string, error) {
	out, errs := f.stream.Start(items, failAt)
	values := make([]string, 0)
	for value := range out {
		values = append(values, value)
	}
	return values, <-errs
}
