package service

import (
	"context"

	"flowcontrol/internal/store"
)

type AlertEvaluationCancelFlow struct {
	state *store.AlertEvaluationCancelState
}

func NewAlertEvaluationCancelFlow(state *store.AlertEvaluationCancelState) *AlertEvaluationCancelFlow {
	return &AlertEvaluationCancelFlow{state: state}
}

func (f *AlertEvaluationCancelFlow) Dispatch(ctx context.Context, key string) error {
	return f.state.Commit(context.Background(), key)
}
