package service

import (
	"sort"
	"time"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/idgen"
)

// recordDecision 写入一条决策日志。
func (s *Service) recordDecision(resource, key string, allowed bool, reason string) {
	d := &model.DecisionLog{
		ID:        idgen.Hex(),
		Resource:  resource,
		Key:       key,
		Allowed:   allowed,
		Reason:    reason,
		CreatedAt: time.Now(),
	}
	_ = s.store.CreateDecisionLog(d)
}

// ListDecisionLogs 按筛选条件返回决策日志（时间倒序）。
func (s *Service) ListDecisionLogs(filter model.DecisionLogFilter, limit int) []*model.DecisionLog {
	if limit <= 0 {
		limit = 100
	}
	all := s.store.ListDecisionLogs()
	matched := make([]*model.DecisionLog, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	if len(matched) > limit {
		matched = matched[:limit]
	}
	return matched
}

// DecisionStats 决策统计。
type DecisionStats struct {
	Total    int     `json:"total"`
	Allowed  int     `json:"allowed"`
	Rejected int     `json:"rejected"`
	AllowRate float64 `json:"allow_rate"`
}

// DecisionStats 返回决策日志统计。
func (s *Service) DecisionStats(resource string) DecisionStats {
	stats := DecisionStats{}
	for _, d := range s.store.ListDecisionLogs() {
		if resource != "" && d.Resource != resource {
			continue
		}
		stats.Total++
		if d.Allowed {
			stats.Allowed++
		} else {
			stats.Rejected++
		}
	}
	if stats.Total > 0 {
		stats.AllowRate = float64(stats.Allowed) / float64(stats.Total)
	}
	return stats
}
