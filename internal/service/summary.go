package service

// SystemSummary 系统总览。
type SystemSummary struct {
	FlowStats     FlowStats           `json:"flow_stats"`
	HealthReport  HealthReport        `json:"health_report"`
	BreakerHealth BreakerHealthSummary `json:"breaker_health"`
	QuotaReport   QuotaReport         `json:"quota_report"`
	AlertAnalysis AlertAnalysis       `json:"alert_analysis"`
	ConfigReport  ConfigReport        `json:"config_report"`
}

// SystemSummary 返回系统总览，聚合各类统计与健康信息。
func (s *Service) SystemSummary() SystemSummary {
	return SystemSummary{
		FlowStats:     s.StatsFlow(),
		HealthReport:  s.HealthReport(),
		BreakerHealth: s.BreakerHealthSummary(),
		QuotaReport:   s.QuotaReport(),
		AlertAnalysis: s.AlertAnalysis(),
		ConfigReport:  s.ConfigReport(),
	}
}

// RuleSummary 规则汇总摘要。
type RuleSummary struct {
	Total       int               `json:"total"`
	Enabled     int               `json:"enabled"`
	Disabled    int               `json:"disabled"`
	ByAlgorithm map[string]int    `json:"by_algorithm"`
}

// RuleSummary 返回限流规则汇总摘要。
func (s *Service) RuleSummary() RuleSummary {
	summary := RuleSummary{ByAlgorithm: make(map[string]int)}
	for _, r := range s.store.ListRateLimitRules() {
		summary.Total++
		if r.Status == "enabled" {
			summary.Enabled++
		} else {
			summary.Disabled++
		}
		summary.ByAlgorithm[r.Algorithm]++
	}
	return summary
}
