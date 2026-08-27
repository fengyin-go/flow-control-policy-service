package service

import (
	"sort"
	"time"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/idgen"
)

// recordTraffic 更新资源流量统计。
func (s *Service) recordTraffic(resource string, allowed bool) {
	stats, err := s.store.GetTrafficStatsByResource(resource)
	if err != nil {
		now := time.Now()
		stats = &model.TrafficStats{
			ID:          idgen.Hex(),
			Resource:    resource,
			WindowStart: now,
			WindowEnd:   now.Add(time.Minute),
			UpdatedAt:   now,
		}
		_ = s.store.CreateTrafficStats(stats)
	}
	stats.TotalRequests++
	if allowed {
		stats.AllowedRequests++
	} else {
		stats.RejectedRequests++
	}
	stats.UpdatedAt = time.Now()
	_ = s.store.UpdateTrafficStats(stats)
}

// ListTrafficStats 返回全部流量统计。
func (s *Service) ListTrafficStats() []*model.TrafficStats {
	list := s.store.ListTrafficStats()
	sort.Slice(list, func(i, j int) bool { return list[i].Resource < list[j].Resource })
	return list
}

// TrafficSummary 资源流量汇总。
type TrafficSummary struct {
	Resource      string  `json:"resource"`
	Total         int64   `json:"total"`
	Allowed       int64   `json:"allowed"`
	Rejected      int64   `json:"rejected"`
	RejectRate    float64 `json:"reject_rate"`
}

// TrafficSummary 返回某资源的流量汇总。
func (s *Service) TrafficSummary(resource string) (*TrafficSummary, error) {
	stats, err := s.store.GetTrafficStatsByResource(resource)
	if err != nil {
		return nil, err
	}
	sum := &TrafficSummary{
		Resource: resource,
		Total:    stats.TotalRequests,
		Allowed:  stats.AllowedRequests,
		Rejected: stats.RejectedRequests,
	}
	if stats.TotalRequests > 0 {
		sum.RejectRate = float64(stats.RejectedRequests) / float64(stats.TotalRequests)
	}
	return sum, nil
}
