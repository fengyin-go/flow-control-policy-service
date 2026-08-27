package store

type AlertConditionFilterResolver interface {
	Resolve(string) string
}

type alert_condition_filterResolver struct{ prefix string }

func (r *alert_condition_filterResolver) Resolve(key string) string {
	return r.prefix + key
}

func NewAlertConditionFilterResolver(enabled bool) AlertConditionFilterResolver {
	var resolver *alert_condition_filterResolver
	if enabled {
		resolver = &alert_condition_filterResolver{prefix: "ready:"}
	}
	return resolver
}
