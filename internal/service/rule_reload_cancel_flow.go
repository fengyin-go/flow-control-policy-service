package service

import (
	"context"

	"flowcontrol/internal/store"
)

type RuleReloadCancelFlow struct{ state *store.RuleReloadCancelState }

func NewRuleReloadCancelFlow(state *store.RuleReloadCancelState) *RuleReloadCancelFlow {
	return &RuleReloadCancelFlow{state: state}
}

func (f *RuleReloadCancelFlow) Dispatch(ctx context.Context, key string) error {
	return f.state.Commit(context.Background(), key)
}
