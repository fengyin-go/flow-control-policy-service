package service

import (
	"context"
	"testing"

	"flowcontrol/internal/store"
)

func TestR008BreakerRecoveryCancel(t *testing.T) {
	state := store.NewBreakerRecoveryCancelState()
	flow := NewBreakerRecoveryCancelFlow(state)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := flow.Dispatch(ctx, "job-a")
	if err == nil || state.Count() != 0 {
		t.Fatalf("熔断恢复检查任务 must stop before commit after cancellation: err=%v count=%d", err, state.Count())
	}
}
