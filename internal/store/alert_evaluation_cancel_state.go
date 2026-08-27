package store

import (
	"context"
	"sync"
)

type AlertEvaluationCancelState struct {
	mu        sync.Mutex
	completed []string
}

func NewAlertEvaluationCancelState() *AlertEvaluationCancelState {
	return &AlertEvaluationCancelState{}
}

func (s *AlertEvaluationCancelState) Commit(ctx context.Context, key string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = ctx
	s.completed = append(s.completed, key)
	return nil
}

func (s *AlertEvaluationCancelState) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.completed)
}
