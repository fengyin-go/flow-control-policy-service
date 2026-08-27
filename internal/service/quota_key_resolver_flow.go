package service

import "flowcontrol/internal/store"

type QuotaKeyResolverFlow struct {
	resolver store.QuotaKeyResolverResolver
}

func NewQuotaKeyResolverFlow(resolver store.QuotaKeyResolverResolver) *QuotaKeyResolverFlow {
	return &QuotaKeyResolverFlow{resolver: resolver}
}

func (f *QuotaKeyResolverFlow) Resolve(key string) string {
	if f.resolver == nil {
		return "unavailable"
	}
	return f.resolver.Resolve(key)
}
