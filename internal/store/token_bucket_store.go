package store

import "flowcontrol/internal/model"

func (s *MemoryStore) CreateTokenBucket(t *model.TokenBucket) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.tokenBuckets {
		if exist.RuleID == t.RuleID {
			return ErrConflict
		}
	}
	s.tokenBuckets[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTokenBucket(id string) (*model.TokenBucket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tokenBuckets[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) GetTokenBucketByRule(ruleID string) (*model.TokenBucket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.tokenBuckets {
		if t.RuleID == ruleID {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListTokenBuckets() []*model.TokenBucket {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TokenBucket, 0, len(s.tokenBuckets))
	for _, t := range s.tokenBuckets {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTokenBucket(t *model.TokenBucket) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tokenBuckets[t.ID]; !ok {
		return ErrNotFound
	}
	s.tokenBuckets[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTokenBucket(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tokenBuckets[id]; !ok {
		return ErrNotFound
	}
	delete(s.tokenBuckets, id)
	return nil
}
