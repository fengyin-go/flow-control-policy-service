package service

import (
	"testing"

	"flowcontrol/internal/store"
)

func TestR024QuotaResetLease(t *testing.T) {
	state := store.NewQuotaResetLeaseLeaseState()
	flow := NewQuotaResetLeaseFlow(state)
	firstErr := flow.Process("node-a", true)
	secondErr := flow.Process("node-a", false)
	if firstErr == nil || secondErr != nil || state.Active() != 0 {
		t.Fatalf("配额重置租约 must be released after failure so the next run can start: first=%v second=%v active=%d", firstErr, secondErr, state.Active())
	}
}
