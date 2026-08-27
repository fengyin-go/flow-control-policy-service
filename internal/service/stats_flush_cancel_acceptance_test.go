package service

import (
	"context"
	"testing"

	"flowcontrol/internal/store"
)

func TestR006StatsFlushCancel(t *testing.T) {
	state := store.NewStatsFlushCancelState()
	flow := NewStatsFlushCancelFlow(state)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := flow.Dispatch(ctx, "job-a")
	if err == nil || state.Count() != 0 {
		t.Fatalf("流量统计刷新任务 must stop before commit after cancellation: err=%v count=%d", err, state.Count())
	}
}
