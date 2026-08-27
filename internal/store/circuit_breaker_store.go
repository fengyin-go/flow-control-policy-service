package store

import "flowcontrol/internal/model"

func (s *MemoryStore) CreateCircuitBreaker(c *model.CircuitBreaker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.breakers {
		if exist.Resource == c.Resource {
			return ErrConflict
		}
	}
	s.breakers[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCircuitBreaker(id string) (*model.CircuitBreaker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.breakers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) GetCircuitBreakerByResource(resource string) (*model.CircuitBreaker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.breakers {
		if c.Resource == resource {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListCircuitBreakers() []*model.CircuitBreaker {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.CircuitBreaker, 0, len(s.breakers))
	for _, c := range s.breakers {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateCircuitBreaker(c *model.CircuitBreaker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.breakers[c.ID]; !ok {
		return ErrNotFound
	}
	s.breakers[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteCircuitBreaker(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.breakers[id]; !ok {
		return ErrNotFound
	}
	delete(s.breakers, id)
	return nil
}
