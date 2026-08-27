package store

import "flowcontrol/internal/model"

func (s *MemoryStore) CreateRateLimitRule(r *model.RateLimitRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.rules {
		if exist.Resource == r.Resource {
			return ErrConflict
		}
	}
	s.rules[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRateLimitRule(id string) (*model.RateLimitRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) GetRateLimitRuleByResource(resource string) (*model.RateLimitRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.rules {
		if r.Resource == resource {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRateLimitRules() []*model.RateLimitRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RateLimitRule, 0, len(s.rules))
	for _, r := range s.rules {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRateLimitRule(r *model.RateLimitRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[r.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.rules {
		if exist.ID != r.ID && exist.Resource == r.Resource {
			return ErrConflict
		}
	}
	s.rules[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRateLimitRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[id]; !ok {
		return ErrNotFound
	}
	delete(s.rules, id)
	return nil
}
