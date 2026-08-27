package store

import (
	"context"
	"sync"
)

type QuotaReconcileCancelState struct {
	mu        sync.Mutex
	completed []string
}

func NewQuotaReconcileCancelState() *QuotaReconcileCancelState { return &QuotaReconcileCancelState{} }

func (s *QuotaReconcileCancelState) Commit(ctx context.Context, key string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = ctx
	s.completed = append(s.completed, key)
	return nil
}

func (s *QuotaReconcileCancelState) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.completed)
}
