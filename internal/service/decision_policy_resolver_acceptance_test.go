package service

import (
	"testing"

	"flowcontrol/internal/store"
)

func TestR011DecisionPolicyResolver(t *testing.T) {
	flow := NewDecisionPolicyResolverFlow(store.NewDecisionPolicyResolverResolver(false))
	panicked := false
	value := ""
	func() {
		defer func() { panicked = recover() != nil }()
		value = flow.Resolve("node-a")
	}()
	if panicked || value != "unavailable" {
		t.Fatalf("放行策略解析组件 must use the disabled fallback without panic: panic=%v value=%q", panicked, value)
	}
}
