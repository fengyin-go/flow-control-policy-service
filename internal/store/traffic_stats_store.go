package store

import "flowcontrol/internal/model"

func (s *MemoryStore) CreateTrafficStats(t *model.TrafficStats) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.trafficStats {
		if exist.Resource == t.Resource {
			return ErrConflict
		}
	}
	s.trafficStats[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTrafficStats(id string) (*model.TrafficStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.trafficStats[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) GetTrafficStatsByResource(resource string) (*model.TrafficStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.trafficStats {
		if t.Resource == resource {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListTrafficStats() []*model.TrafficStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TrafficStats, 0, len(s.trafficStats))
	for _, t := range s.trafficStats {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTrafficStats(t *model.TrafficStats) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.trafficStats[t.ID]; !ok {
		return ErrNotFound
	}
	s.trafficStats[t.ID] = t
	return nil
}
