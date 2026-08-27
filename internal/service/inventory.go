package service

import (
	"sort"

	"flowcontrol/internal/model"
)

// ResourceInventory 资源清单条目。
type ResourceInventory struct {
	Resource   string                `json:"resource"`
	Rule       *model.RateLimitRule  `json:"rule,omitempty"`
	Breaker    *model.CircuitBreaker `json:"breaker,omitempty"`
	AlertCount int                   `json:"alert_count"`
	QuotaCount int                   `json:"quota_count"`
	HasStats   bool                  `json:"has_stats"`
}

// ResourceInventory 返回全部资源的完整清单（规则 + 熔断器 + 告警 + 配额 + 统计）。
func (s *Service) ResourceInventory() []ResourceInventory {
	resources := make(map[string]*ResourceInventory)
	ensure := func(res string) *ResourceInventory {
		inv, ok := resources[res]
		if !ok {
			inv = &ResourceInventory{Resource: res}
			resources[res] = inv
		}
		return inv
	}
	for _, r := range s.store.ListRateLimitRules() {
		ensure(r.Resource).Rule = r
	}
	for _, b := range s.store.ListCircuitBreakers() {
		ensure(b.Resource).Breaker = b
	}
	for _, a := range s.store.ListAlertRules() {
		if a.Resource != "" {
			ensure(a.Resource).AlertCount++
		}
	}
	for _, q := range s.store.ListQuotas() {
		if rule, err := s.store.GetRateLimitRule(q.RuleID); err == nil {
			ensure(rule.Resource).QuotaCount++
		}
	}
	for _, t := range s.store.ListTrafficStats() {
		ensure(t.Resource).HasStats = true
	}
	out := make([]ResourceInventory, 0, len(resources))
	for _, inv := range resources {
		out = append(out, *inv)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Resource < out[j].Resource })
	return out
}
