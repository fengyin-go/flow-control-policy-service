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
	// 无论本轮重置成功还是失败，都必须结束租约，否则重试同一维度时会被残留的忙碌状态拒绝。
	defer f.state.Release(key, token)
	if fail {
		return errors.New("operation failed")
	}
	return nil
}
