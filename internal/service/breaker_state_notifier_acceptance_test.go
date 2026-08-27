package service

import (
	"testing"

	"flowcontrol/internal/store"
)

func TestR013BreakerStateNotifier(t *testing.T) {
	flow := NewBreakerStateNotifierFlow(store.NewBreakerStateNotifierResolver(false))
	panicked := false
	value := ""
	func() {
		defer func() { panicked = recover() != nil }()
		value = flow.Resolve("node-a")
	}()
	if panicked || value != "unavailable" {
		t.Fatalf("熔断状态通知组件 must use the disabled fallback without panic: panic=%v value=%q", panicked, value)
	}
}
