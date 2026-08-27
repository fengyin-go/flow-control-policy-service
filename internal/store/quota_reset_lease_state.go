package store

import "sync"

type QuotaResetLeaseLeaseState struct {
	mu     sync.Mutex
	active map[string]uint64
	next   uint64
}

func NewQuotaResetLeaseLeaseState() *QuotaResetLeaseLeaseState {
	return &QuotaResetLeaseLeaseState{active: make(map[string]uint64)}
}

func (s *QuotaResetLeaseLeaseState) Acquire(key string) (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.active[key]; exists {
		return 0, false
	}
	s.next++
	token := s.next
	s.active[key] = token
	return token, true
}

func (s *QuotaResetLeaseLeaseState) Release(key string, token uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.active[key]
	if !ok || current != token {
		// 租约已被其它轮次占用或已释放，拒绝释放不属于本轮的资源状态。
		return false
	}
	delete(s.active, key)
	return true
}

func (s *QuotaResetLeaseLeaseState) Active() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.active)
}
