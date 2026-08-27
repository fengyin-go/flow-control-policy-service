package service

import (
	"testing"

	"flowcontrol/internal/store"
)

func TestR022BucketRefillLease(t *testing.T) {
	state := store.NewBucketRefillLeaseLeaseState()
	flow := NewBucketRefillLeaseFlow(state)
	firstErr := flow.Process("node-a", true)
	secondErr := flow.Process("node-a", false)
	if firstErr == nil || secondErr != nil || state.Active() != 0 {
		t.Fatalf("令牌桶补充租约 must be released after failure so the next run can start: first=%v second=%v active=%d", firstErr, secondErr, state.Active())
	}
}
