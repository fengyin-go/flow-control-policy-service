package service

import (
	"context"

	"flowcontrol/internal/store"
)

type QuotaReconcileCancelFlow struct {
	state *store.QuotaReconcileCancelState
}

func NewQuotaReconcileCancelFlow(state *store.QuotaReconcileCancelState) *QuotaReconcileCancelFlow {
	return &QuotaReconcileCancelFlow{state: state}
}

func (f *QuotaReconcileCancelFlow) Dispatch(ctx context.Context, key string) error {
	return f.state.Commit(context.Background(), key)
}
