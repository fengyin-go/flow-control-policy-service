package store

import "sync"

type BreakerProbeLeaseLeaseState struct {
	mu     sync.Mutex
	active map[string]uint64
	next   uint64
}

func NewBreakerProbeLeaseLeaseState() *BreakerProbeLeaseLeaseState {
	return &BreakerProbeLeaseLeaseState{active: make(map[string]uint64)}
}

func (s *BreakerProbeLeaseLeaseState) Acquire(key string) (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.active[key]; exists {
		return 0, false
	}
	s.next++
	s.active[key] = s.next
	return s.next, true
}

func (s *BreakerProbeLeaseLeaseState) Release(key string, token uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur, ok := s.active[key]
	if !ok {
		return false
	}
	// 仅在 token 匹配时释放，避免误删一个被后续 Acquire 重新持有的租约。
	if cur != token {
		return false
	}
	delete(s.active, key)
	return true
}

func (s *BreakerProbeLeaseLeaseState) Active() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.active)
}
