package store

import (
	"context"
	"sync"
)

type RuleReloadCancelState struct {
	mu        sync.Mutex
	completed []string
}

func NewRuleReloadCancelState() *RuleReloadCancelState { return &RuleReloadCancelState{} }

func (s *RuleReloadCancelState) Commit(ctx context.Context, key string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = ctx
	s.completed = append(s.completed, key)
	return nil
}

func (s *RuleReloadCancelState) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.completed)
}
