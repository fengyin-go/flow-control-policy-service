package service

import "flowcontrol/internal/store"

type BreakerStateNotifierFlow struct {
	resolver store.BreakerStateNotifierResolver
}

func NewBreakerStateNotifierFlow(resolver store.BreakerStateNotifierResolver) *BreakerStateNotifierFlow {
	return &BreakerStateNotifierFlow{resolver: resolver}
}

func (f *BreakerStateNotifierFlow) Resolve(key string) string {
	if f.resolver == nil {
		return "unavailable"
	}
	return f.resolver.Resolve(key)
}
