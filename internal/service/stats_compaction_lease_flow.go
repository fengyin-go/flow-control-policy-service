package service

import (
	"errors"

	"flowcontrol/internal/store"
)

type StatsCompactionLeaseFlow struct {
	state *store.StatsCompactionLeaseLeaseState
}

func NewStatsCompactionLeaseFlow(state *store.StatsCompactionLeaseLeaseState) *StatsCompactionLeaseFlow {
	return &StatsCompactionLeaseFlow{state: state}
}

func (f *StatsCompactionLeaseFlow) Process(key string, fail bool) error {
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
