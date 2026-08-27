package service

import "flowcontrol/internal/store"

type DecisionBatchFlow struct{ stream *store.DecisionBatchStream }

func NewDecisionBatchFlow(stream *store.DecisionBatchStream) *DecisionBatchFlow {
	return &DecisionBatchFlow{stream: stream}
}

func (f *DecisionBatchFlow) Collect(items []string, failAt int) ([]string, error) {
	out, errs := f.stream.Start(items, failAt)
	values := make([]string, 0)
	for value := range out {
		values = append(values, value)
	}
	return values, <-errs
}
