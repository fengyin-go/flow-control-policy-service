package service

import "flowcontrol/internal/store"

type AlertDeliveryBatchFlow struct {
	stream *store.AlertDeliveryBatchStream
}

func NewAlertDeliveryBatchFlow(stream *store.AlertDeliveryBatchStream) *AlertDeliveryBatchFlow {
	return &AlertDeliveryBatchFlow{stream: stream}
}

// Collect consumes the delivery stream and returns whatever results arrived
// before the stream ended. If the stream reports an error, collection stops
// immediately and the partial results collected so far are returned together
// with that error.
func (f *AlertDeliveryBatchFlow) Collect(items []string, failAt int) ([]string, error) {
	out, errs := f.stream.Start(items, failAt)
	values := make([]string, 0)
	for {
		select {
		case value, ok := <-out:
			if !ok {
				// out closed: the goroutine has returned, so errs is already
				// closed (defer order) and this receive is non-blocking. It
				// yields the error when the stream was interrupted, or nil
				// when it completed normally.
				return values, <-errs
			}
			values = append(values, value)
		case err := <-errs:
			// Error signaled: stop immediately, drain any out value still
			// pending until out is closed, then return partial results with
			// the error.
			for value := range out {
				values = append(values, value)
			}
			return values, err
		}
	}
}
