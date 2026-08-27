package service

import "flowcontrol/internal/store"

type AlertStateRetryFlow struct {
	state *store.AlertStateRetryRetryState
}

func NewAlertStateRetryFlow(state *store.AlertStateRetryRetryState) *AlertStateRetryFlow {
	return &AlertStateRetryFlow{state: state}
}

func (f *AlertStateRetryFlow) Execute() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		last = f.state.Next()
		if last == nil {
			return nil
		}
	}
	return last
}
