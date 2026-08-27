package service

import (
	"testing"

	"flowcontrol/internal/store"
)

func TestR029AlertStateRetry(t *testing.T) {
	permanent := store.NewAlertStateRetryRetryState(&store.RetryFailure{Temporary: false, Message: "rejected"}, nil)
	permanentErr := NewAlertStateRetryFlow(permanent).Execute()
	temporary := store.NewAlertStateRetryRetryState(&store.RetryFailure{Temporary: true, Message: "busy"}, nil)
	temporaryErr := NewAlertStateRetryFlow(temporary).Execute()
	if permanentErr == nil || permanent.Attempts() != 1 || permanent.Commits() != 0 || temporaryErr != nil || temporary.Attempts() != 2 || temporary.Commits() != 1 {
		t.Fatalf("告警状态写入 must stop permanent failures and retry temporary failures exactly once: permanent=(%v,%d,%d) temporary=(%v,%d,%d)", permanentErr, permanent.Attempts(), permanent.Commits(), temporaryErr, temporary.Attempts(), temporary.Commits())
	}
}
