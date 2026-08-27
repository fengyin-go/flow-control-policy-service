package service

import (
	"sort"

	"flowcontrol/internal/model"
)

// ThresholdAdvice 限流阈值建议。
type ThresholdAdvice struct {
	Resource        string `json:"resource"`
	CurrentLimit    int    `json:"current_limit"`
	PeakRequests    int64  `json:"peak_requests"`
	SuggestedLimit  int    `json:"suggested_limit"`
	Action          string `json:"action"`
}

// AdviseThresholds 基于流量统计给出限流阈值建议。
func (s *Service) AdviseThresholds() []ThresholdAdvice {
	advice := make([]ThresholdAdvice, 0)
	for _, t := range s.store.ListTrafficStats() {
		item := ThresholdAdvice{
			Resource:     t.Resource,
			PeakRequests: t.TotalRequests,
		}
		if rule, err := s.store.GetRateLimitRuleByResource(t.Resource); err == nil {
			item.CurrentLimit = rule.Limit
		}
		suggested := int(float64(t.TotalRequests) * 1.2)
		if suggested < 10 {
			suggested = 10
		}
		item.SuggestedLimit = suggested
		if item.CurrentLimit == 0 {
			item.Action = "create"
		} else if suggested > item.CurrentLimit {
			item.Action = "increase"
		} else if suggested < item.CurrentLimit/2 {
			item.Action = "decrease"
		} else {
			item.Action = "keep"
		}
		advice = append(advice, item)
	}
	sort.Slice(advice, func(i, j int) bool { return advice[i].PeakRequests > advice[j].PeakRequests })
	return advice
}

// RuleStatusReport 规则状态报告。
func (s *Service) RuleStatusReport() map[string]int {
	return map[string]int{
		"total":    len(s.store.ListRateLimitRules()),
		"enabled":  countRulesByStatus(s.store.ListRateLimitRules(), model.RuleEnabled),
		"disabled": countRulesByStatus(s.store.ListRateLimitRules(), model.RuleDisabled),
	}
}

func countRulesByStatus(rules []*model.RateLimitRule, status string) int {
	n := 0
	for _, r := range rules {
		if r.Status == status {
			n++
		}
	}
	return n
}
