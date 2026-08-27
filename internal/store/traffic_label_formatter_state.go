package store

type TrafficLabelFormatterResolver interface {
	Resolve(string) string
}

type traffic_label_formatterResolver struct{ prefix string }

func (r *traffic_label_formatterResolver) Resolve(key string) string {
	return r.prefix + key
}

func NewTrafficLabelFormatterResolver(enabled bool) TrafficLabelFormatterResolver {
	var resolver *traffic_label_formatterResolver
	if enabled {
		resolver = &traffic_label_formatterResolver{prefix: "ready:"}
	}
	return resolver
}
