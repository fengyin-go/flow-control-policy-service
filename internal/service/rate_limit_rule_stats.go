package service

import (
	"sort"

	"flowcontrol/internal/model"
)

// RuleStats 单规则统计。
type RuleStats struct {
	Rule          *model.RateLimitRule `json:"rule"`
	Total         int64                `json:"total"`
	Rejected      int64                `json:"rejected"`
	RejectRate    float64              `json:"reject_rate"`
	DecisionCount int                  `json:"decision_count"`
}

// RuleStatsReport 返回全部规则的统计（按拒绝率降序）。
func (s *Service) RuleStatsReport() []RuleStats {
	out := make([]RuleStats, 0)
	for _, r := range s.store.ListRateLimitRules() {
		rs := RuleStats{Rule: r}
		if t, err := s.store.GetTrafficStatsByResource(r.Resource); err == nil {
			rs.Total = t.TotalRequests
			rs.Rejected = t.RejectedRequests
			if t.TotalRequests > 0 {
				rs.RejectRate = float64(t.RejectedRequests) / float64(t.TotalRequests)
			}
		}
		rs.DecisionCount = s.DecisionStats(r.Resource).Total
		out = append(out, rs)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].RejectRate > out[j].RejectRate })
	return out
}

// DisabledRuleResources 返回被停用规则保护的资源。
func (s *Service) DisabledRuleResources() []string {
	out := make([]string, 0)
	for _, r := range s.store.ListRateLimitRules() {
		if r.Status == model.RuleDisabled {
			out = append(out, r.Resource)
		}
	}
	sort.Strings(out)
	return out
}

// UnmatchedBreakers 返回没有对应流量统计的熔断器资源。
func (s *Service) UnmatchedBreakers() []string {
	out := make([]string, 0)
	for _, b := range s.store.ListCircuitBreakers() {
		if _, err := s.store.GetTrafficStatsByResource(b.Resource); err != nil {
			out = append(out, b.Resource)
		}
	}
	sort.Strings(out)
	return out
}
