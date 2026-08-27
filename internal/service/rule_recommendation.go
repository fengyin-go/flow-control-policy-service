package service

import (
	"fmt"

	"flowcontrol/internal/model"
)

// RuleRecommendation 规则建议。
type RuleRecommendation struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Detail   string `json:"detail"`
}

// RuleRecommendations 返回基于流量统计的规则调整建议。
func (s *Service) RuleRecommendations() []RuleRecommendation {
	recs := make([]RuleRecommendation, 0)
	for _, a := range s.AdviseThresholds() {
		recs = append(recs, RuleRecommendation{
			Resource: a.Resource,
			Action:   a.Action,
			Detail:   fmt.Sprintf("当前限流 %d，建议 %d", a.CurrentLimit, a.SuggestedLimit),
		})
	}
	return recs
}

// QuotaAnalysis 配额分析。
type QuotaAnalysis struct {
	Total       int            `json:"total"`
	Exhausted   int            `json:"exhausted"`
	HighUsage   int            `json:"high_usage"`
	ByDimension map[string]int `json:"by_dimension"`
}

// QuotaAnalysis 分析配额使用情况。
func (s *Service) QuotaAnalysis() QuotaAnalysis {
	a := QuotaAnalysis{ByDimension: make(map[string]int)}
	for _, q := range s.store.ListQuotas() {
		a.Total++
		a.ByDimension[q.Dimension]++
		if q.Used >= q.Allowed {
			a.Exhausted++
		} else if q.Allowed > 0 && float64(q.Used)/float64(q.Allowed) > 0.8 {
			a.HighUsage++
		}
	}
	return a
}

// ConfigReport 配置汇总。
type ConfigReport struct {
	DefaultLimit     int `json:"default_limit"`
	DefaultWindowSec int `json:"default_window_sec"`
	EnabledRules     int `json:"enabled_rules"`
	DisabledRules    int `json:"disabled_rules"`
	OpenBreakers     int `json:"open_breakers"`
	EnabledAlerts    int `json:"enabled_alerts"`
}

// ConfigReport 返回运行配置与规则状态汇总。
func (s *Service) ConfigReport() ConfigReport {
	r := ConfigReport{}
	if s.cfg != nil {
		r.DefaultLimit = s.cfg.DefaultLimit
		r.DefaultWindowSec = s.cfg.DefaultWindowSec
	}
	for _, rule := range s.store.ListRateLimitRules() {
		if rule.Status == model.RuleEnabled {
			r.EnabledRules++
		} else {
			r.DisabledRules++
		}
	}
	for _, b := range s.store.ListCircuitBreakers() {
		if b.State == model.BreakerOpen {
			r.OpenBreakers++
		}
	}
	for _, a := range s.store.ListAlertRules() {
		if a.Status == model.RuleEnabled {
			r.EnabledAlerts++
		}
	}
	return r
}
