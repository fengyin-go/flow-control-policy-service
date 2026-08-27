package service

import (
	"testing"
	"time"

	"flowcontrol/internal/store"
)

func TestR017QuotaConsumeBatch(t *testing.T) {
	flow := NewQuotaConsumeBatchFlow(store.NewQuotaConsumeBatchStream())
	items := []string{"node-a", "node-b", "node-c"}
	type outcome struct {
		values []string
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		values, err := flow.Collect(items, 1)
		done <- outcome{values: values, err: err}
	}()
	got := outcome{}
	timedOut := false
	select {
	case got = <-done:
	case <-time.After(500 * time.Millisecond):
		timedOut = true
	}
	if timedOut || got.err == nil || len(got.values) != 1 {
		t.Fatalf("配额消耗批次 must close promptly and preserve the partial result on an error: timeout=%v err=%v values=%v", timedOut, got.err, got.values)
	}
}
