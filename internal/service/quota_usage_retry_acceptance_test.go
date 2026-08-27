package service

import (
	"testing"

	"flowcontrol/internal/store"
)

func TestR030QuotaUsageRetry(t *testing.T) {
	permanent := store.NewQuotaUsageRetryRetryState(&store.RetryFailure{Temporary: false, Message: "rejected"}, nil)
	permanentErr := NewQuotaUsageRetryFlow(permanent).Execute()
	temporary := store.NewQuotaUsageRetryRetryState(&store.RetryFailure{Temporary: true, Message: "busy"}, nil)
	temporaryErr := NewQuotaUsageRetryFlow(temporary).Execute()
	if permanentErr == nil || permanent.Attempts() != 1 || permanent.Commits() != 0 || temporaryErr != nil || temporary.Attempts() != 2 || temporary.Commits() != 1 {
		t.Fatalf("配额用量写入 must stop permanent failures and retry temporary failures exactly once: permanent=(%v,%d,%d) temporary=(%v,%d,%d)", permanentErr, permanent.Attempts(), permanent.Commits(), temporaryErr, temporary.Attempts(), temporary.Commits())
	}
}
