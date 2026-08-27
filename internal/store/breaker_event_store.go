package store

import "flowcontrol/internal/model"

func (s *MemoryStore) CreateBreakerEvent(e *model.BreakerEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.breakerEvents[e.ID] = e
	return nil
}

func (s *MemoryStore) ListBreakerEvents() []*model.BreakerEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.BreakerEvent, 0, len(s.breakerEvents))
	for _, e := range s.breakerEvents {
		list = append(list, e)
	}
	return list
}
