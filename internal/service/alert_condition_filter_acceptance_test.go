package service

import (
	"testing"

	"flowcontrol/internal/store"
)

func TestR014AlertConditionFilter(t *testing.T) {
	flow := NewAlertConditionFilterFlow(store.NewAlertConditionFilterResolver(false))
	panicked := false
	value := ""
	func() {
		defer func() { panicked = recover() != nil }()
		value = flow.Resolve("node-a")
	}()
	if panicked || value != "unavailable" {
		t.Fatalf("告警条件筛选组件 must use the disabled fallback without panic: panic=%v value=%q", panicked, value)
	}
}
