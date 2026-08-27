package service

import (
	"testing"

	"flowcontrol/internal/store"
)

func TestR025StatsCompactionLease(t *testing.T) {
	state := store.NewStatsCompactionLeaseLeaseState()
	flow := NewStatsCompactionLeaseFlow(state)
	firstErr := flow.Process("node-a", true)
	secondErr := flow.Process("node-a", false)
	if firstErr == nil || secondErr != nil || state.Active() != 0 {
		t.Fatalf("统计压缩租约 must be released after failure so the next run can start: first=%v second=%v active=%d", firstErr, secondErr, state.Active())
	}
}
