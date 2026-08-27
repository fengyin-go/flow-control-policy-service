package service

import (
	"context"
	"testing"

	"flowcontrol/internal/store"
)

func TestR010RuleReloadCancel(t *testing.T) {
	state := store.NewRuleReloadCancelState()
	flow := NewRuleReloadCancelFlow(state)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := flow.Dispatch(ctx, "job-a")
	if err == nil || state.Count() != 0 {
		t.Fatalf("限流规则重载任务 must stop before commit after cancellation: err=%v count=%d", err, state.Count())
	}
}
