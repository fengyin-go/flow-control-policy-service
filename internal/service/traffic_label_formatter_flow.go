package service

import "flowcontrol/internal/store"

type TrafficLabelFormatterFlow struct {
	resolver store.TrafficLabelFormatterResolver
}

func NewTrafficLabelFormatterFlow(resolver store.TrafficLabelFormatterResolver) *TrafficLabelFormatterFlow {
	return &TrafficLabelFormatterFlow{resolver: resolver}
}

func (f *TrafficLabelFormatterFlow) Resolve(key string) string {
	if f.resolver == nil {
		return "unavailable"
	}
	return f.resolver.Resolve(key)
}
