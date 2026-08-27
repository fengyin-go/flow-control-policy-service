package service

import (
	"errors"

	"flowcontrol/internal/store"
)

type RuleRefreshLeaseFlow struct {
	state *store.RuleRefreshLeaseLeaseState
}

func NewRuleRefreshLeaseFlow(state *store.RuleRefreshLeaseLeaseState) *RuleRefreshLeaseFlow {
	return &RuleRefreshLeaseFlow{state: state}
}

func (f *RuleRefreshLeaseFlow) Process(key string, fail bool) error {
	token, ok := f.state.Acquire(key)
	if !ok {
		return errors.New("lease busy")
	}
	if fail {
		return errors.New("operation failed")
	}
	defer f.state.Release(key, token)
	return nil
}
