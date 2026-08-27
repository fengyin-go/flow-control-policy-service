package store

type BreakerStateNotifierResolver interface {
	Resolve(string) string
}

type breaker_state_notifierResolver struct{ prefix string }

func (r *breaker_state_notifierResolver) Resolve(key string) string {
	return r.prefix + key
}

func NewBreakerStateNotifierResolver(enabled bool) BreakerStateNotifierResolver {
	var resolver *breaker_state_notifierResolver
	if enabled {
		resolver = &breaker_state_notifierResolver{prefix: "ready:"}
	}
	return resolver
}
