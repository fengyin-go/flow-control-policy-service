package service

import (
	"errors"

	"flowcontrol/internal/store"
)

type BreakerEventRetryFlow struct {
	state *store.BreakerEventRetryRetryState
}

func NewBreakerEventRetryFlow(state *store.BreakerEventRetryRetryState) *BreakerEventRetryFlow {
	return &BreakerEventRetryFlow{state: state}
}

// Execute 逐步执行熔断事件写入，仅对临时错误重试。
// 遇到永久错误（非临时）立即返回，不再重试，避免同一事件重复写入并提交。
func (f *BreakerEventRetryFlow) Execute() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		last = f.state.Next()
		if last == nil {
			return nil
		}
		if !isTemporary(last) {
			return last
		}
	}
	return last
}

// isTemporary 判断错误是否属于可重试的临时错误。
// 仅承认 store.RetryFailure 且 Temporary 为真的情况，其余一律视为永久错误。
func isTemporary(err error) bool {
	var rf *store.RetryFailure
	if as := errors.As(err, &rf); as {
		return rf.Temporary
	}
	return false
}
