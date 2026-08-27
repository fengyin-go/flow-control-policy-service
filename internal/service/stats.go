package service

import (
	"flowcontrol/internal/model"
)

// FlowStats 流量控制整体统计。
type FlowStats struct {
	TotalRules        int            `json:"total_rules"`
	TotalBreakers     int            `json:"total_breakers"`
	TotalQuotas       int            `json:"total_quotas"`
	TotalAlertRules   int            `json:"total_alert_rules"`
	BreakerByState    map[string]int `json:"breaker_by_state"`
	TotalRequests     int64          `json:"total_requests"`
	RejectedRequests  int64          `json:"rejected_requests"`
	OverallRejectRate float64        `json:"overall_reject_rate"`
}

// StatsFlow 返回流量控制整体统计。
func (s *Service) StatsFlow() FlowStats {
	stats := FlowStats{BreakerByState: make(map[string]int)}
	stats.TotalRules = len(s.store.ListRateLimitRules())
	stats.TotalBreakers = len(s.store.ListCircuitBreakers())
	stats.TotalQuotas = len(s.store.ListQuotas())
	stats.TotalAlertRules = len(s.store.ListAlertRules())
	for _, b := range s.store.ListCircuitBreakers() {
		stats.BreakerByState[b.State]++
	}
	for _, t := range s.store.ListTrafficStats() {
		stats.TotalRequests += t.TotalRequests
		stats.RejectedRequests += t.RejectedRequests
	}
	if stats.TotalRequests > 0 {
		stats.OverallRejectRate = float64(stats.RejectedRequests) / float64(stats.TotalRequests)
	}
	return stats
}

// ResourceStats 单资源统计。
type ResourceStats struct {
	Resource   string                 `json:"resource"`
	Rule       *model.RateLimitRule   `json:"rule,omitempty"`
	Breaker    *model.CircuitBreaker  `json:"breaker,omitempty"`
	Total      int64                  `json:"total"`
	Rejected   int64                  `json:"rejected"`
	RejectRate float64                `json:"reject_rate"`
}

// ResourceStats 返回某资源的限流与熔断综合统计。
func (s *Service) ResourceStats(resource string) ResourceStats {
	rs := ResourceStats{Resource: resource}
	if rule, err := s.store.GetRateLimitRuleByResource(resource); err == nil {
		rs.Rule = rule
	}
	if b, err := s.store.GetCircuitBreakerByResource(resource); err == nil {
		rs.Breaker = b
	}
	if t, err := s.store.GetTrafficStatsByResource(resource); err == nil {
		rs.Total = t.TotalRequests
		rs.Rejected = t.RejectedRequests
		if t.TotalRequests > 0 {
			rs.RejectRate = float64(t.RejectedRequests) / float64(t.TotalRequests)
		}
	}
	return rs
}
