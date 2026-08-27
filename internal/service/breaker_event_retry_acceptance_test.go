package service

import (
	"testing"

	"flowcontrol/internal/store"
)

func TestR027BreakerEventRetry(t *testing.T) {
	permanent := store.NewBreakerEventRetryRetryState(&store.RetryFailure{Temporary: false, Message: "rejected"}, nil)
	permanentErr := NewBreakerEventRetryFlow(permanent).Execute()
	temporary := store.NewBreakerEventRetryRetryState(&store.RetryFailure{Temporary: true, Message: "busy"}, nil)
	temporaryErr := NewBreakerEventRetryFlow(temporary).Execute()
	if permanentErr == nil || permanent.Attempts() != 1 || permanent.Commits() != 0 || temporaryErr != nil || temporary.Attempts() != 2 || temporary.Commits() != 1 {
		t.Fatalf("熔断事件写入 must stop permanent failures and retry temporary failures exactly once: permanent=(%v,%d,%d) temporary=(%v,%d,%d)", permanentErr, permanent.Attempts(), permanent.Commits(), temporaryErr, temporary.Attempts(), temporary.Commits())
	}
}
