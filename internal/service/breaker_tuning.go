package service

import (
	"flowcontrol/internal/model"
)

// BreakerTuningSuggestion 熔断器调优建议。
type BreakerTuningSuggestion struct {
	Breaker    *model.CircuitBreaker `json:"breaker"`
	Suggestion string                `json:"suggestion"`
	Reason     string                `json:"reason"`
}

// BreakerTuningSuggestions 返回熔断器调优建议。
func (s *Service) BreakerTuningSuggestions() []BreakerTuningSuggestion {
	out := make([]BreakerTuningSuggestion, 0)
	for _, b := range s.store.ListCircuitBreakers() {
		sug := ""
		reason := ""
		openEvents := 0
		for _, e := range s.store.ListBreakerEvents() {
			if e.BreakerID == b.ID && e.EventType == model.BreakerEventOpened {
				openEvents++
			}
		}
		switch {
		case openEvents >= 3:
			sug = "提高失败阈值或延长半开观察期"
			reason = "熔断打开过于频繁"
		case b.FailureThreshold < 3:
			sug = "适当提高失败阈值，避免误熔断"
			reason = "失败阈值过低"
		case b.OpenDurationSec < 10:
			sug = "延长熔断打开时长，给下游更多恢复时间"
			reason = "打开时长过短"
		}
		if sug != "" {
			out = append(out, BreakerTuningSuggestion{Breaker: b, Suggestion: sug, Reason: reason})
		}
	}
	return out
}

// BreakerEffectiveness 熔断器有效性。
type BreakerEffectiveness struct {
	Breaker          *model.CircuitBreaker `json:"breaker"`
	OpenEvents       int                   `json:"open_events"`
	RejectedRequests int64                 `json:"rejected_requests"`
	Effective        bool                  `json:"effective"`
}

// BreakerEffectiveness 评估熔断器是否真正拦截了流量。
func (s *Service) BreakerEffectiveness() []BreakerEffectiveness {
	out := make([]BreakerEffectiveness, 0)
	for _, b := range s.store.ListCircuitBreakers() {
		eff := BreakerEffectiveness{Breaker: b}
		for _, e := range s.store.ListBreakerEvents() {
			if e.BreakerID == b.ID && e.EventType == model.BreakerEventOpened {
				eff.OpenEvents++
			}
		}
		if t, err := s.store.GetTrafficStatsByResource(b.Resource); err == nil {
			eff.RejectedRequests = t.RejectedRequests
		}
		eff.Effective = eff.OpenEvents > 0 || eff.RejectedRequests > 0
		out = append(out, eff)
	}
	return out
}
