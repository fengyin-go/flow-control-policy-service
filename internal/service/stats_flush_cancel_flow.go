package service

import (
	"context"

	"flowcontrol/internal/store"
)

type StatsFlushCancelFlow struct{ state *store.StatsFlushCancelState }

func NewStatsFlushCancelFlow(state *store.StatsFlushCancelState) *StatsFlushCancelFlow {
	return &StatsFlushCancelFlow{state: state}
}

func (f *StatsFlushCancelFlow) Dispatch(ctx context.Context, key string) error {
	return f.state.Commit(context.Background(), key)
}
