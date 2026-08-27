package store

import (
	"context"
	"sync"
)

type StatsFlushCancelState struct {
	mu        sync.Mutex
	completed []string
}

func NewStatsFlushCancelState() *StatsFlushCancelState { return &StatsFlushCancelState{} }

func (s *StatsFlushCancelState) Commit(ctx context.Context, key string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = ctx
	s.completed = append(s.completed, key)
	return nil
}

func (s *StatsFlushCancelState) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.completed)
}
