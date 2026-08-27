package store

type QuotaKeyResolverResolver interface {
	Resolve(string) string
}

type quota_key_resolverResolver struct{ prefix string }

func (r *quota_key_resolverResolver) Resolve(key string) string {
	return r.prefix + key
}

func NewQuotaKeyResolverResolver(enabled bool) QuotaKeyResolverResolver {
	var resolver *quota_key_resolverResolver
	if enabled {
		resolver = &quota_key_resolverResolver{prefix: "ready:"}
	}
	return resolver
}
