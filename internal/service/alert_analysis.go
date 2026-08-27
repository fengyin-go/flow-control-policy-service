package service

import (
	"sort"
)

// AlertAnalysis 告警分析。
type AlertAnalysis struct {
	TotalRules    int            `json:"total_rules"`
	EnabledRules  int            `json:"enabled_rules"`
	TriggeredHits int            `json:"triggered_hits"`
	ByMetric      map[string]int `json:"by_metric"`
	Hits          []AlertHit     `json:"hits"`
}

// AlertAnalysis 分析告警规则与当前命中情况。
func (s *Service) AlertAnalysis() AlertAnalysis {
	analysis := AlertAnalysis{ByMetric: make(map[string]int)}
	for _, r := range s.store.ListAlertRules() {
		analysis.TotalRules++
		if r.Status == "enabled" {
			analysis.EnabledRules++
		}
		analysis.ByMetric[r.Metric]++
	}
	analysis.Hits = s.EvaluateAlerts()
	analysis.TriggeredHits = len(analysis.Hits)
	return analysis
}

// AlertRuleRank 告警规则触发排行。
type AlertRuleRank struct {
	Name  string `json:"name"`
	Value float64 `json:"value"`
}

// TopTriggeredAlerts 返回触发值最高的告警规则。
func (s *Service) TopTriggeredAlerts(n int) []AlertRuleRank {
	if n <= 0 {
		n = 10
	}
	hits := s.EvaluateAlerts()
	list := make([]AlertRuleRank, 0, len(hits))
	for _, h := range hits {
		list = append(list, AlertRuleRank{Name: h.Rule.Name, Value: h.MetricValue})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Value > list[j].Value })
	if len(list) > n {
		list = list[:n]
	}
	return list
}
