package store

import (
	"context"
	"sync"
)

type BreakerRecoveryCancelState struct {
	mu        sync.Mutex
	completed []string
}

func NewBreakerRecoveryCancelState() *BreakerRecoveryCancelState {
	return &BreakerRecoveryCancelState{}
}

func (s *BreakerRecoveryCancelState) Commit(ctx context.Context, key string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = ctx
	s.completed = append(s.completed, key)
	return nil
}

func (s *BreakerRecoveryCancelState) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.completed)
}
