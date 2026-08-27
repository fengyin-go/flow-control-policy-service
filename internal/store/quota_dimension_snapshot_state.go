package store

import "sync"

type QuotaDimensionSnapshotState struct {
	mu     sync.RWMutex
	values map[string][]string
}

func NewQuotaDimensionSnapshotState() *QuotaDimensionSnapshotState {
	return &QuotaDimensionSnapshotState{values: make(map[string][]string)}
}

func (s *QuotaDimensionSnapshotState) Replace(key string, values []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = values
}

func (s *QuotaDimensionSnapshotState) Snapshot(key string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.values[key]
}
