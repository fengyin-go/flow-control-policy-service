package service

import (
	"context"

	"flowcontrol/internal/store"
)

type BreakerRecoveryCancelFlow struct {
	state *store.BreakerRecoveryCancelState
}

func NewBreakerRecoveryCancelFlow(state *store.BreakerRecoveryCancelState) *BreakerRecoveryCancelFlow {
	return &BreakerRecoveryCancelFlow{state: state}
}

func (f *BreakerRecoveryCancelFlow) Dispatch(ctx context.Context, key string) error {
	return f.state.Commit(context.Background(), key)
}
