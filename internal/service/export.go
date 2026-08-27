package service

import (
	"flowcontrol/internal/model"
	"flowcontrol/pkg/idgen"
)

// Snapshot 流量控制全量快照。
type Snapshot struct {
	Rules         []*model.RateLimitRule  `json:"rules"`
	TokenBuckets  []*model.TokenBucket    `json:"token_buckets"`
	Breakers      []*model.CircuitBreaker `json:"breakers"`
	BreakerEvents []*model.BreakerEvent   `json:"breaker_events"`
	Quotas        []*model.Quota          `json:"quotas"`
	TrafficStats  []*model.TrafficStats   `json:"traffic_stats"`
	AlertRules    []*model.AlertRule      `json:"alert_rules"`
}

// ExportSnapshot 导出全量快照。
func (s *Service) ExportSnapshot() Snapshot {
	return Snapshot{
		Rules:         s.store.ListRateLimitRules(),
		TokenBuckets:  s.store.ListTokenBuckets(),
		Breakers:      s.store.ListCircuitBreakers(),
		BreakerEvents: s.store.ListBreakerEvents(),
		Quotas:        s.store.ListQuotas(),
		TrafficStats:  s.store.ListTrafficStats(),
		AlertRules:    s.store.ListAlertRules(),
	}
}

// ImportSnapshot 导入快照（跳过重复项），返回各实体导入数量。
func (s *Service) ImportSnapshot(snap Snapshot) (map[string]int, error) {
	imported := map[string]int{
		"rules":          0,
		"token_buckets":  0,
		"breakers":       0,
		"breaker_events": 0,
		"quotas":         0,
		"traffic_stats":  0,
		"alert_rules":    0,
	}
	for _, r := range snap.Rules {
		if r.ID == "" {
			r.ID = idgen.Hex()
		}
		if err := s.store.CreateRateLimitRule(r); err == nil {
			imported["rules"]++
		}
	}
	for _, t := range snap.TokenBuckets {
		if t.ID == "" {
			t.ID = idgen.Hex()
		}
		if err := s.store.CreateTokenBucket(t); err == nil {
			imported["token_buckets"]++
		}
	}
	for _, b := range snap.Breakers {
		if b.ID == "" {
			b.ID = idgen.Hex()
		}
		if err := s.store.CreateCircuitBreaker(b); err == nil {
			imported["breakers"]++
		}
	}
	for _, e := range snap.BreakerEvents {
		if e.ID == "" {
			e.ID = idgen.Hex()
		}
		if err := s.store.CreateBreakerEvent(e); err == nil {
			imported["breaker_events"]++
		}
	}
	for _, q := range snap.Quotas {
		if q.ID == "" {
			q.ID = idgen.Hex()
		}
		if err := s.store.CreateQuota(q); err == nil {
			imported["quotas"]++
		}
	}
	for _, t := range snap.TrafficStats {
		if t.ID == "" {
			t.ID = idgen.Hex()
		}
		if err := s.store.CreateTrafficStats(t); err == nil {
			imported["traffic_stats"]++
		}
	}
	for _, a := range snap.AlertRules {
		if a.ID == "" {
			a.ID = idgen.Hex()
		}
		if err := s.store.CreateAlertRule(a); err == nil {
			imported["alert_rules"]++
		}
	}
	return imported, nil
}
