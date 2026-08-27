package service

import "flowcontrol/internal/store"

type DecisionLogRetryFlow struct {
	state *store.DecisionLogRetryRetryState
}

func NewDecisionLogRetryFlow(state *store.DecisionLogRetryRetryState) *DecisionLogRetryFlow {
	return &DecisionLogRetryFlow{state: state}
}

func (f *DecisionLogRetryFlow) Execute() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		last = f.state.Next()
		if last == nil {
			return nil
		}
	}
	return last
}
