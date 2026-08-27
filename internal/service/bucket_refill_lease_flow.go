package service

import (
	"errors"

	"flowcontrol/internal/store"
)

type BucketRefillLeaseFlow struct {
	state *store.BucketRefillLeaseLeaseState
}

func NewBucketRefillLeaseFlow(state *store.BucketRefillLeaseLeaseState) *BucketRefillLeaseFlow {
	return &BucketRefillLeaseFlow{state: state}
}

func (f *BucketRefillLeaseFlow) Process(key string, fail bool) error {
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
