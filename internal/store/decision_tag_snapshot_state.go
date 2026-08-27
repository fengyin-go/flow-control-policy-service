package store

import "sync"

type DecisionTagSnapshotState struct {
	mu     sync.RWMutex
	values map[string][]string
}

func NewDecisionTagSnapshotState() *DecisionTagSnapshotState {
	return &DecisionTagSnapshotState{values: make(map[string][]string)}
}

func (s *DecisionTagSnapshotState) Replace(key string, values []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = values
}

func (s *DecisionTagSnapshotState) Snapshot(key string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.values[key]
}
