package service

import (
	"sort"
)

// DecisionTrend 决策趋势（按资源聚合）。
type DecisionTrend struct {
	Resource    string  `json:"resource"`
	Total       int     `json:"total"`
	Allowed     int     `json:"allowed"`
	Rejected    int     `json:"rejected"`
	RejectRate  float64 `json:"reject_rate"`
}

// DecisionTrends 返回按资源聚合的决策趋势。
func (s *Service) DecisionTrends() []DecisionTrend {
	agg := make(map[string]*DecisionTrend)
	for _, d := range s.store.ListDecisionLogs() {
		t, ok := agg[d.Resource]
		if !ok {
			t = &DecisionTrend{Resource: d.Resource}
			agg[d.Resource] = t
		}
		t.Total++
		if d.Allowed {
			t.Allowed++
		} else {
			t.Rejected++
		}
	}
	out := make([]DecisionTrend, 0, len(agg))
	for _, t := range agg {
		if t.Total > 0 {
			t.RejectRate = float64(t.Rejected) / float64(t.Total)
		}
		out = append(out, *t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].RejectRate > out[j].RejectRate })
	return out
}

// TopRejectedKeys 返回被拒绝次数最多的 N 个维度键。
type KeyRejection struct {
	Key      string `json:"key"`
	Rejected int    `json:"rejected"`
}

// TopRejectedKeys 返回被拒绝次数最多的维度键。
func (s *Service) TopRejectedKeys(n int) []KeyRejection {
	if n <= 0 {
		n = 10
	}
	counts := make(map[string]int)
	for _, d := range s.store.ListDecisionLogs() {
		if !d.Allowed {
			counts[d.Key]++
		}
	}
	list := make([]KeyRejection, 0, len(counts))
	for k, v := range counts {
		list = append(list, KeyRejection{Key: k, Rejected: v})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Rejected > list[j].Rejected })
	if len(list) > n {
		list = list[:n]
	}
	return list
}

// ReasonDistribution 拒绝原因分布。
func (s *Service) ReasonDistribution() map[string]int {
	dist := make(map[string]int)
	for _, d := range s.store.ListDecisionLogs() {
		if !d.Allowed {
			dist[d.Reason]++
		}
	}
	return dist
}
