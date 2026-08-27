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
	// 统计压缩无论成功还是失败，返回前都必须清理活动租约；
	// 否则同一 key 的下一次压缩会一直看到资源占用。
	defer f.state.Release(key, token)
	if fail {
		return errors.New("operation failed")
	}
	return nil
}
