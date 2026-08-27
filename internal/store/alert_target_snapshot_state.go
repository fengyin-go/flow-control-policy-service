package store

import "sync"

type AlertTargetSnapshotState struct {
	mu     sync.RWMutex
	values map[string][]string
}

func NewAlertTargetSnapshotState() *AlertTargetSnapshotState {
	return &AlertTargetSnapshotState{values: make(map[string][]string)}
}

func (s *AlertTargetSnapshotState) Replace(key string, values []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = values
}

func (s *AlertTargetSnapshotState) Snapshot(key string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.values[key]
}
