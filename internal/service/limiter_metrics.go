package service

import (
	"time"

	"flowcontrol/internal/model"
)

// LimiterRuntimeMetrics 限流器运行时指标。
type LimiterRuntimeMetrics struct {
	ActiveWindows int `json:"active_windows"`
	SlidingKeys   int `json:"sliding_keys"`
	BucketKeys    int `json:"bucket_keys"`
}

// LimiterRuntimeMetrics 返回限流器运行时状态规模。
func (s *Service) LimiterRuntimeMetrics() LimiterRuntimeMetrics {
	s.lim.mu.Lock()
	defer s.lim.mu.Unlock()
	return LimiterRuntimeMetrics{
		ActiveWindows: len(s.lim.windows),
		SlidingKeys:   len(s.lim.sliding),
		BucketKeys:    len(s.lim.buckets),
	}
}

// ResetLimiter 清空限流器运行时状态。
func (s *Service) ResetLimiter() {
	s.lim.mu.Lock()
	defer s.lim.mu.Unlock()
	s.lim.windows = make(map[string]*windowState)
	s.lim.sliding = make(map[string][]time.Time)
	s.lim.buckets = make(map[string]*bucketState)
}

// RuleRuntimeView 规则运行时视图。
type RuleRuntimeView struct {
	Rule      *model.RateLimitRule `json:"rule"`
	Algorithm string               `json:"algorithm"`
	Limit     int                  `json:"limit"`
	WindowSec int                  `json:"window_sec"`
}

// RuleRuntimeViews 返回全部启用规则的运行时视图。
func (s *Service) RuleRuntimeViews() []RuleRuntimeView {
	views := make([]RuleRuntimeView, 0)
	for _, r := range s.store.ListRateLimitRules() {
		if r.Status != model.RuleEnabled {
			continue
		}
		views = append(views, RuleRuntimeView{
			Rule:      r,
			Algorithm: r.Algorithm,
			Limit:     r.Limit,
			WindowSec: r.WindowSec,
		})
	}
	return views
}
