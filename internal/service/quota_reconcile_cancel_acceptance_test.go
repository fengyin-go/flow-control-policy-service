package service

import (
	"context"
	"testing"

	"flowcontrol/internal/store"
)

func TestR007QuotaReconcileCancel(t *testing.T) {
	state := store.NewQuotaReconcileCancelState()
	flow := NewQuotaReconcileCancelFlow(state)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := flow.Dispatch(ctx, "job-a")
	if err == nil || state.Count() != 0 {
		t.Fatalf("配额校准任务 must stop before commit after cancellation: err=%v count=%d", err, state.Count())
	}
}
