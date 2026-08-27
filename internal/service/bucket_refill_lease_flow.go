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
	// 拿到租约后立即注册释放，确保后续失败路径也会归还占用标记，
	// 否则残留租约会挡住下一次对同一个桶的补充。
	defer f.state.Release(key, token)
	if fail {
		return errors.New("operation failed")
	}
	return nil
}
