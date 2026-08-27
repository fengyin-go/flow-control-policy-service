package service

import (
	"time"

	"flowcontrol/internal/model"
)

// BreakerEventTimeline 熔断事件时间线节点。
type BreakerEventTimeline struct {
	EventType string    `json:"event_type"`
	At        time.Time `json:"at"`
	Detail    string    `json:"detail"`
}

// BreakerEventTimeline 返回某熔断器的事件时间线（时间升序）。
func (s *Service) BreakerEventTimeline(breakerID string) []BreakerEventTimeline {
	events := s.ListBreakerEvents(breakerID)
	out := make([]BreakerEventTimeline, 0, len(events))
	// ListBreakerEvents 为时间倒序，反转得到升序
	for i := len(events) - 1; i >= 0; i-- {
		e := events[i]
		out = append(out, BreakerEventTimeline{EventType: e.EventType, At: e.CreatedAt, Detail: e.Detail})
	}
	return out
}

// EventSequenceAnalysis 熔断事件序列分析。
type EventSequenceAnalysis struct {
	Opened    int  `json:"opened"`
	Closed    int  `json:"closed"`
	HalfOpen  int  `json:"half_open"`
	Rejected  int  `json:"rejected"`
	Flapping  bool `json:"flapping"`
}

// EventSequenceAnalysis 分析某熔断器的事件序列。
func (s *Service) EventSequenceAnalysis(breakerID string) EventSequenceAnalysis {
	analysis := EventSequenceAnalysis{}
	for _, e := range s.store.ListBreakerEvents() {
		if e.BreakerID != breakerID {
			continue
		}
		switch e.EventType {
		case model.BreakerEventOpened:
			analysis.Opened++
		case model.BreakerEventClosed:
			analysis.Closed++
		case model.BreakerEventHalfOpen:
			analysis.HalfOpen++
		case model.BreakerEventRejected:
			analysis.Rejected++
		}
	}
	analysis.Flapping = analysis.Opened >= 3
	return analysis
}
