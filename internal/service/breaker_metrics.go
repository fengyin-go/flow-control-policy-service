package service

import (
	"sort"
	"time"

	"flowcontrol/internal/model"
)

// BreakerMetrics 熔断器指标。
type BreakerMetrics struct {
	Breaker       *model.CircuitBreaker `json:"breaker"`
	OpenEvents    int                   `json:"open_events"`
	CloseEvents   int                   `json:"close_events"`
	HalfOpenEvents int                  `json:"half_open_events"`
	RejectEvents  int                   `json:"reject_events"`
	StateAgeSeconds int64               `json:"state_age_seconds"`
}

// BreakerMetrics 返回某熔断器的指标。
func (s *Service) BreakerMetrics(breakerID string) (*BreakerMetrics, error) {
	b, err := s.store.GetCircuitBreaker(breakerID)
	if err != nil {
		return nil, err
	}
	m := &BreakerMetrics{
		Breaker:        b,
		StateAgeSeconds: int64(time.Since(b.LastStateChange).Seconds()),
	}
	for _, e := range s.store.ListBreakerEvents() {
		if e.BreakerID != breakerID {
			continue
		}
		switch e.EventType {
		case model.BreakerEventOpened:
			m.OpenEvents++
		case model.BreakerEventClosed:
			m.CloseEvents++
		case model.BreakerEventHalfOpen:
			m.HalfOpenEvents++
		case model.BreakerEventRejected:
			m.RejectEvents++
		}
	}
	return m, nil
}

// BreakerRank 熔断器活跃度排行项。
type BreakerRank struct {
	Breaker *model.CircuitBreaker `json:"breaker"`
	Events  int                   `json:"events"`
}

// MostActiveBreakers 按事件数对熔断器排行。
func (s *Service) MostActiveBreakers(n int) []BreakerRank {
	if n <= 0 {
		n = 10
	}
	counts := make(map[string]int)
	for _, e := range s.store.ListBreakerEvents() {
		counts[e.BreakerID]++
	}
	list := make([]BreakerRank, 0)
	for _, b := range s.store.ListCircuitBreakers() {
		list = append(list, BreakerRank{Breaker: b, Events: counts[b.ID]})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Events > list[j].Events })
	if len(list) > n {
		list = list[:n]
	}
	return list
}
