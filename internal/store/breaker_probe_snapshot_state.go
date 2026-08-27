package store

import "sync"

type BreakerProbeSnapshotState struct {
	mu     sync.RWMutex
	values map[string][]string
}

func NewBreakerProbeSnapshotState() *BreakerProbeSnapshotState {
	return &BreakerProbeSnapshotState{values: make(map[string][]string)}
}

func (s *BreakerProbeSnapshotState) Replace(key string, values []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = values
}

func (s *BreakerProbeSnapshotState) Snapshot(key string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.values[key]
}
