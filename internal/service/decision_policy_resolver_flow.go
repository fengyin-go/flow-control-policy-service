package service

import "flowcontrol/internal/store"

type DecisionPolicyResolverFlow struct {
	resolver store.DecisionPolicyResolverResolver
}

func NewDecisionPolicyResolverFlow(resolver store.DecisionPolicyResolverResolver) *DecisionPolicyResolverFlow {
	return &DecisionPolicyResolverFlow{resolver: resolver}
}

func (f *DecisionPolicyResolverFlow) Resolve(key string) string {
	if f.resolver == nil {
		return "unavailable"
	}
	return f.resolver.Resolve(key)
}
