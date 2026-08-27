package store

type DecisionPolicyResolverResolver interface {
	Resolve(string) string
}

type decision_policy_resolverResolver struct{ prefix string }

func (r *decision_policy_resolverResolver) Resolve(key string) string {
	return r.prefix + key
}

func NewDecisionPolicyResolverResolver(enabled bool) DecisionPolicyResolverResolver {
	var resolver *decision_policy_resolverResolver
	if enabled {
		resolver = &decision_policy_resolverResolver{prefix: "ready:"}
	}
	return resolver
}
