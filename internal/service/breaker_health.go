package service

import (
	"time"

	"flowcontrol/internal/model"
)

// BreakerHealthSummary 熔断器健康汇总。
type BreakerHealthSummary struct {
	Total    int            `json:"total"`
	ByState  map[string]int `json:"by_state"`
	OpenBreakers []string   `json:"open_breakers"`
	HealthGrade  string     `json:"health_grade"`
}

// BreakerHealthSummary 返回熔断器健康汇总。
func (s *Service) BreakerHealthSummary() BreakerHealthSummary {
	sum := BreakerHealthSummary{ByState: make(map[string]int), OpenBreakers: []string{}}
	for _, b := range s.store.ListCircuitBreakers() {
		sum.Total++
		sum.ByState[b.State]++
		if b.State == model.BreakerOpen {
			sum.OpenBreakers = append(sum.OpenBreakers, b.Resource)
		}
	}
	switch {
	case sum.ByState[model.BreakerOpen] == 0:
		sum.HealthGrade = "healthy"
	case sum.ByState[model.BreakerOpen] < sum.Total:
		sum.HealthGrade = "degraded"
	default:
		sum.HealthGrade = "critical"
	}
	return sum
}

// BreakerRecoveryTime 熔断器预计恢复时间。
type BreakerRecovery struct {
	Breaker       *model.CircuitBreaker `json:"breaker"`
	RemainingSec  int64                 `json:"remaining_seconds"`
	Recovered     bool                  `json:"recovered"`
}

// BreakerRecoveries 返回所有 open 状态熔断器的预计恢复时间。
func (s *Service) BreakerRecoveries() []BreakerRecovery {
	out := make([]BreakerRecovery, 0)
	now := time.Now()
	for _, b := range s.store.ListCircuitBreakers() {
		if b.State != model.BreakerOpen {
			continue
		}
		elapsed := now.Sub(b.LastStateChange).Seconds()
		need := float64(b.OpenDurationSec)
		rec := BreakerRecovery{Breaker: b, Recovered: elapsed >= need}
		if elapsed < need {
			rec.RemainingSec = int64(need - elapsed)
		}
		out = append(out, rec)
	}
	return out
}
