package service

import "flowcontrol/internal/store"

type AlertConditionFilterFlow struct {
	resolver store.AlertConditionFilterResolver
}

func NewAlertConditionFilterFlow(resolver store.AlertConditionFilterResolver) *AlertConditionFilterFlow {
	return &AlertConditionFilterFlow{resolver: resolver}
}

func (f *AlertConditionFilterFlow) Resolve(key string) string {
	if f.resolver == nil {
		return "unavailable"
	}
	return f.resolver.Resolve(key)
}
