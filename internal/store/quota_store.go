package store

import "flowcontrol/internal/model"

func (s *MemoryStore) CreateQuota(q *model.Quota) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.quotas {
		if exist.RuleID == q.RuleID && exist.Dimension == q.Dimension && exist.DimensionValue == q.DimensionValue {
			return ErrConflict
		}
	}
	s.quotas[q.ID] = q
	return nil
}

func (s *MemoryStore) GetQuota(id string) (*model.Quota, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.quotas[id]
	if !ok {
		return nil, ErrNotFound
	}
	return q, nil
}

func (s *MemoryStore) ListQuotas() []*model.Quota {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Quota, 0, len(s.quotas))
	for _, q := range s.quotas {
		list = append(list, q)
	}
	return list
}

func (s *MemoryStore) UpdateQuota(q *model.Quota) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.quotas[q.ID]; !ok {
		return ErrNotFound
	}
	s.quotas[q.ID] = q
	return nil
}

func (s *MemoryStore) DeleteQuota(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.quotas[id]; !ok {
		return ErrNotFound
	}
	delete(s.quotas, id)
	return nil
}
