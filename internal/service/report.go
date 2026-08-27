package service

import (
	"sort"

	"flowcontrol/internal/model"
)

// RuleDistribution 限流算法分布。
type RuleDistribution struct {
	ByAlgorithm map[string]int `json:"by_algorithm"`
	ByStatus    map[string]int `json:"by_status"`
}

// RuleDistribution 返回限流规则的算法与状态分布。
func (s *Service) RuleDistribution() RuleDistribution {
	d := RuleDistribution{
		ByAlgorithm: make(map[string]int),
		ByStatus:    make(map[string]int),
	}
	for _, r := range s.store.ListRateLimitRules() {
		d.ByAlgorithm[r.Algorithm]++
		d.ByStatus[r.Status]++
	}
	return d
}

// RejectedResource 拒绝量排行项。
type RejectedResource struct {
	Resource string  `json:"resource"`
	Rejected int64   `json:"rejected"`
	Total    int64   `json:"total"`
	Rate     float64 `json:"rate"`
}

// TopRejectedResources 返回拒绝率最高的 N 个资源。
func (s *Service) TopRejectedResources(n int) []RejectedResource {
	if n <= 0 {
		n = 10
	}
	list := make([]RejectedResource, 0)
	for _, t := range s.store.ListTrafficStats() {
		item := RejectedResource{
			Resource: t.Resource,
			Rejected: t.RejectedRequests,
			Total:    t.TotalRequests,
		}
		if t.TotalRequests > 0 {
			item.Rate = float64(t.RejectedRequests) / float64(t.TotalRequests)
		}
		list = append(list, item)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Rate > list[j].Rate })
	if len(list) > n {
		list = list[:n]
	}
	return list
}

// BreakerEventDistribution 熔断事件类型分布。
func (s *Service) BreakerEventDistribution() map[string]int {
	dist := make(map[string]int)
	for _, e := range s.store.ListBreakerEvents() {
		dist[e.EventType]++
	}
	return dist
}

// QuotaUsage 配额使用率排行项。
type QuotaUsage struct {
	Quota     *model.Quota `json:"quota"`
	UsageRate float64      `json:"usage_rate"`
}

// TopQuotaUsage 返回配额使用率最高的 N 个配额。
func (s *Service) TopQuotaUsage(n int) []QuotaUsage {
	if n <= 0 {
		n = 10
	}
	list := make([]QuotaUsage, 0)
	for _, q := range s.store.ListQuotas() {
		u := QuotaUsage{Quota: q}
		if q.Allowed > 0 {
			u.UsageRate = float64(q.Used) / float64(q.Allowed)
		}
		list = append(list, u)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].UsageRate > list[j].UsageRate })
	if len(list) > n {
		list = list[:n]
	}
	return list
}
