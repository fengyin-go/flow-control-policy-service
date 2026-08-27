package service

import (
	"sort"

	"flowcontrol/internal/model"
)

// OptimizationReport 限流优化报告。
type OptimizationReport struct {
	UnprotectedResources []string `json:"unprotected_resources"`
	OverProvisionedRules []string `json:"over_provisioned_rules"`
}

// OptimizationReport 找出无规则保护但已有流量的资源，以及限流过松的规则。
func (s *Service) OptimizationReport() OptimizationReport {
	report := OptimizationReport{
		UnprotectedResources: []string{},
		OverProvisionedRules: []string{},
	}
	for _, t := range s.store.ListTrafficStats() {
		if _, err := s.store.GetRateLimitRuleByResource(t.Resource); err != nil {
			report.UnprotectedResources = append(report.UnprotectedResources, t.Resource)
		}
	}
	for _, r := range s.store.ListRateLimitRules() {
		if stats, err := s.store.GetTrafficStatsByResource(r.Resource); err == nil {
			if stats.TotalRequests > 0 && int64(r.Limit) > stats.TotalRequests*10 {
				report.OverProvisionedRules = append(report.OverProvisionedRules, r.Resource)
			}
		}
	}
	sort.Strings(report.UnprotectedResources)
	sort.Strings(report.OverProvisionedRules)
	return report
}

// RuleComplexityReport 规则复杂度报告。
type RuleComplexityReport struct {
	TotalRules    int     `json:"total_rules"`
	FixedWindow   int     `json:"fixed_window"`
	SlidingWindow int     `json:"sliding_window"`
	TokenBucket   int     `json:"token_bucket"`
	AvgLimit      float64 `json:"avg_limit"`
	MinLimit      int     `json:"min_limit"`
	MaxLimit      int     `json:"max_limit"`
}

// RuleComplexityReport 统计规则算法分布与阈值范围。
func (s *Service) RuleComplexityReport() RuleComplexityReport {
	r := RuleComplexityReport{}
	total := 0
	sum := 0
	for _, rule := range s.store.ListRateLimitRules() {
		total++
		sum += rule.Limit
		switch rule.Algorithm {
		case model.AlgorithmFixedWindow:
			r.FixedWindow++
		case model.AlgorithmSlidingWindow:
			r.SlidingWindow++
		case model.AlgorithmTokenBucket:
			r.TokenBucket++
		}
		if r.MinLimit == 0 || rule.Limit < r.MinLimit {
			r.MinLimit = rule.Limit
		}
		if rule.Limit > r.MaxLimit {
			r.MaxLimit = rule.Limit
		}
	}
	r.TotalRules = total
	if total > 0 {
		r.AvgLimit = float64(sum) / float64(total)
	}
	return r
}
