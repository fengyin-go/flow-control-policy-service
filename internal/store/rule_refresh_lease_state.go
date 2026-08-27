package store

import "sync"

type RuleRefreshLeaseLeaseState struct {
	mu     sync.Mutex
	active map[string]uint64
	next   uint64
}

func NewRuleRefreshLeaseLeaseState() *RuleRefreshLeaseLeaseState {
	return &RuleRefreshLeaseLeaseState{active: make(map[string]uint64)}
}

func (s *RuleRefreshLeaseLeaseState) Acquire(key string) (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.active[key]; exists {
		return 0, false
	}
	token := uint64(1)
	s.active[key] = token
	return token, true
}

func (s *RuleRefreshLeaseLeaseState) Release(key string, token uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = token
	if _, ok := s.active[key]; !ok {
		return false
	}
	delete(s.active, key)
	return true
}

func (s *RuleRefreshLeaseLeaseState) Active() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.active)
}
