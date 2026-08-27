package service

import (
	"context"
	"testing"

	"flowcontrol/internal/store"
)

func TestR009AlertEvaluationCancel(t *testing.T) {
	state := store.NewAlertEvaluationCancelState()
	flow := NewAlertEvaluationCancelFlow(state)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := flow.Dispatch(ctx, "job-a")
	if err == nil || state.Count() != 0 {
		t.Fatalf("告警评估任务 must stop before commit after cancellation: err=%v count=%d", err, state.Count())
	}
}
