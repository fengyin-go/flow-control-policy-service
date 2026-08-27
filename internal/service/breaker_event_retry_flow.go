package service

import "flowcontrol/internal/store"

type BreakerEventRetryFlow struct {
	state *store.BreakerEventRetryRetryState
}

func NewBreakerEventRetryFlow(state *store.BreakerEventRetryRetryState) *BreakerEventRetryFlow {
	return &BreakerEventRetryFlow{state: state}
}

func (f *BreakerEventRetryFlow) Execute() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		last = f.state.Next()
		if last == nil {
			return nil
		}
	}
	return last
}
