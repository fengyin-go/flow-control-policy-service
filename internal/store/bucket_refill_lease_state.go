package store

import "sync"

type BucketRefillLeaseLeaseState struct {
	mu     sync.Mutex
	active map[string]uint64
	next   uint64
}

func NewBucketRefillLeaseLeaseState() *BucketRefillLeaseLeaseState {
	return &BucketRefillLeaseLeaseState{active: make(map[string]uint64)}
}

func (s *BucketRefillLeaseLeaseState) Acquire(key string) (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.active[key]; exists {
		return 0, false
	}
	// token 从 0 开始单调递增，0 保留为未携带 token 的哨兵值，
	// 以便 Release 在不校验 token 时仍能归还（向后兼容旧调用）。
	s.next++
	token := s.next
	s.active[key] = token
	return token, true
}

// Release 归还指定 key 的占用标记。仅当 token 与当前持有者匹配时才删除，
// 避免释放已被重新获取的新租约。token 被忽略的旧实现会把值置 0 但保留 key，
// 导致 Acquire 误判该 key 仍被占用，挡住后续补充。
func (s *BucketRefillLeaseLeaseState) Release(key string, token uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur, ok := s.active[key]
	if !ok {
		return false
	}
	if token != 0 && cur != token {
		return false
	}
	delete(s.active, key)
	return true
}

func (s *BucketRefillLeaseLeaseState) Active() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.active)
}
