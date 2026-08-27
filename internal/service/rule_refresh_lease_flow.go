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
	// Register the release before any failing branch: every exit path must
	// release the lease, otherwise a failed refresh leaves the lease held and
	// the next refresh on the same rule reports "lease busy" forever.
	defer f.state.Release(key, token)
	if fail {
		return errors.New("operation failed")
	}
	return nil
}
