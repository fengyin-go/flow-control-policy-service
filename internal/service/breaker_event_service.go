package service

import (
	"sort"
	"time"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/idgen"
)

// emitBreakerEvent 写入一条熔断事件。
func (s *Service) emitBreakerEvent(breakerID, eventType, detail string) {
	e := &model.BreakerEvent{
		ID:        idgen.Hex(),
		BreakerID: breakerID,
		EventType: eventType,
		Detail:    detail,
		CreatedAt: time.Now(),
	}
	_ = s.store.CreateBreakerEvent(e)
}

// ListBreakerEvents 返回熔断事件（可按熔断器筛选，时间倒序）。
func (s *Service) ListBreakerEvents(breakerID string) []*model.BreakerEvent {
	list := make([]*model.BreakerEvent, 0)
	for _, e := range s.store.ListBreakerEvents() {
		if breakerID == "" || e.BreakerID == breakerID {
			list = append(list, e)
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
	return list
}
