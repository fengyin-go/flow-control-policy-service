package service

import (
	"testing"

	"flowcontrol/internal/store"
)

func TestR012QuotaKeyResolver(t *testing.T) {
	flow := NewQuotaKeyResolverFlow(store.NewQuotaKeyResolverResolver(false))
	panicked := false
	value := ""
	func() {
		defer func() { panicked = recover() != nil }()
		value = flow.Resolve("node-a")
	}()
	if panicked || value != "unavailable" {
		t.Fatalf("配额键解析组件 must use the disabled fallback without panic: panic=%v value=%q", panicked, value)
	}
}
