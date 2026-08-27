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
	// 无论探测成功还是失败，都要释放租约，否则下一次对同一熔断器的
	// 探测会因 key 仍被占用而返回 "lease busy"。
	defer f.state.Release(key, token)
	if fail {
		return errors.New("operation failed")
	}
	return nil
}
