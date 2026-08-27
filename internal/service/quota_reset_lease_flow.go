package service

import (
	"errors"

	"flowcontrol/internal/store"
)

type QuotaResetLeaseFlow struct {
	state *store.QuotaResetLeaseLeaseState
}

func NewQuotaResetLeaseFlow(state *store.QuotaResetLeaseLeaseState) *QuotaResetLeaseFlow {
	return &QuotaResetLeaseFlow{state: state}
}

func (f *QuotaResetLeaseFlow) Process(key string, fail bool) error {
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
