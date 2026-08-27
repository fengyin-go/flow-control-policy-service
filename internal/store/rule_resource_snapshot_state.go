package store

import "sync"

type RuleResourceSnapshotState struct {
	mu     sync.RWMutex
	values map[string][]string
}

func NewRuleResourceSnapshotState() *RuleResourceSnapshotState {
	return &RuleResourceSnapshotState{values: make(map[string][]string)}
}

func (s *RuleResourceSnapshotState) Replace(key string, values []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = values
}

func (s *RuleResourceSnapshotState) Snapshot(key string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.values[key]
}
