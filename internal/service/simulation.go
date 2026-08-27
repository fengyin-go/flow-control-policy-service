package service

import (
	"flowcontrol/internal/model"
)

// SimulationResult 流量模拟结果。
type SimulationResult struct {
	Total    int     `json:"total"`
	Allowed  int     `json:"allowed"`
	Rejected int     `json:"rejected"`
	AllowRate float64 `json:"allow_rate"`
}

// SimulateTraffic 对资源模拟 n 次请求，返回放行统计。
func (s *Service) SimulateTraffic(resource, key string, n int) SimulationResult {
	res := SimulationResult{Total: n}
	for i := 0; i < n; i++ {
		if s.Decide(resource, key).Allowed {
			res.Allowed++
		} else {
			res.Rejected++
		}
	}
	if res.Total > 0 {
		res.AllowRate = float64(res.Allowed) / float64(res.Total)
	}
	return res
}

// ResourceStatus 资源综合状态。
type ResourceStatus struct {
	Resource    string                `json:"resource"`
	HasRule     bool                  `json:"has_rule"`
	HasBreaker  bool                  `json:"has_breaker"`
	Rule        *model.RateLimitRule  `json:"rule,omitempty"`
	Breaker     *model.CircuitBreaker `json:"breaker,omitempty"`
	RejectRate  float64               `json:"reject_rate"`
}

// ResourceStatuses 返回全部资源的综合状态（规则 ∪ 熔断器 ∪ 流量统计）。
func (s *Service) ResourceStatuses() []ResourceStatus {
	resources := make(map[string]*ResourceStatus)
	ensure := func(resource string) *ResourceStatus {
		rs, ok := resources[resource]
		if !ok {
			rs = &ResourceStatus{Resource: resource}
			resources[resource] = rs
		}
		return rs
	}
	for _, r := range s.store.ListRateLimitRules() {
		rs := ensure(r.Resource)
		rs.HasRule = true
		rs.Rule = r
	}
	for _, b := range s.store.ListCircuitBreakers() {
		rs := ensure(b.Resource)
		rs.HasBreaker = true
		rs.Breaker = b
	}
	for _, t := range s.store.ListTrafficStats() {
		rs := ensure(t.Resource)
		if t.TotalRequests > 0 {
			rs.RejectRate = float64(t.RejectedRequests) / float64(t.TotalRequests)
		}
	}
	out := make([]ResourceStatus, 0, len(resources))
	for _, rs := range resources {
		out = append(out, *rs)
	}
	return out
}
