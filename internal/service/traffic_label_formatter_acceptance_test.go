package service

import (
	"testing"

	"flowcontrol/internal/store"
)

func TestR015TrafficLabelFormatter(t *testing.T) {
	flow := NewTrafficLabelFormatterFlow(store.NewTrafficLabelFormatterResolver(false))
	panicked := false
	value := ""
	func() {
		defer func() { panicked = recover() != nil }()
		value = flow.Resolve("node-a")
	}()
	if panicked || value != "unavailable" {
		t.Fatalf("流量标签格式组件 must use the disabled fallback without panic: panic=%v value=%q", panicked, value)
	}
}
