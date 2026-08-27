package service

import (
	"flowcontrol/internal/model"
)

// HealthReport 流量控制整体健康报告。
type HealthReport struct {
	Grade             string   `json:"grade"`
	Score             float64  `json:"score"`
	OverallRejectRate float64  `json:"overall_reject_rate"`
	OpenBreakerCount  int      `json:"open_breaker_count"`
	ExhaustedQuotas   int      `json:"exhausted_quotas"`
	AlertHits         int      `json:"alert_hits"`
	Suggestions       []string `json:"suggestions"`
}

// HealthReport 综合拒绝率、熔断器、配额与告警计算健康分与建议。
func (s *Service) HealthReport() HealthReport {
	stats := s.StatsFlow()
	breakerSum := s.BreakerHealthSummary()
	quota := s.QuotaReport()
	alerts := s.EvaluateAlerts()
	report := HealthReport{
		OverallRejectRate: stats.OverallRejectRate,
		OpenBreakerCount:  breakerSum.ByState[model.BreakerOpen],
		ExhaustedQuotas:   quota.ExhaustedQuotas,
		AlertHits:         len(alerts),
		Suggestions:       []string{},
	}
	score := 100.0
	if stats.OverallRejectRate > 0.3 {
		score -= 30
		report.Suggestions = append(report.Suggestions, "整体拒绝率过高，建议调高限流阈值")
	}
	if report.OpenBreakerCount > 0 {
		score -= 20
		report.Suggestions = append(report.Suggestions, "存在打开状态的熔断器，请检查下游服务")
	}
	if quota.ExhaustedQuotas > 0 {
		score -= 10
		report.Suggestions = append(report.Suggestions, "存在耗尽的配额")
	}
	if len(alerts) > 0 {
		score -= 10
		report.Suggestions = append(report.Suggestions, "有告警规则命中")
	}
	if score < 0 {
		score = 0
	}
	report.Score = score
	switch {
	case score >= 90:
		report.Grade = "A"
	case score >= 75:
		report.Grade = "B"
	case score >= 60:
		report.Grade = "C"
	default:
		report.Grade = "D"
	}
	return report
}
