package service

import (
	"errors"

	"flowcontrol/internal/store"
)

type BreakerProbeLeaseFlow struct {
	state *store.BreakerProbeLeaseLeaseState
}

func NewBreakerProbeLeaseFlow(state *store.BreakerProbeLeaseLeaseState) *BreakerProbeLeaseFlow {
	return &BreakerProbeLeaseFlow{state: state}
}

func (f *BreakerProbeLeaseFlow) Process(key string, fail bool) error {
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
