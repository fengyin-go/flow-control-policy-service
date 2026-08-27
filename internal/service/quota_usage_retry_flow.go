package service

import (
	"errors"

	"flowcontrol/internal/store"
)

type QuotaUsageRetryFlow struct {
	state *store.QuotaUsageRetryRetryState
}

func NewQuotaUsageRetryFlow(state *store.QuotaUsageRetryRetryState) *QuotaUsageRetryFlow {
	return &QuotaUsageRetryFlow{state: state}
}

// Execute 执行配额用量写入。仅对临时占用错误重试；
// 不可恢复（永久）拒绝原样返回，不重试，更不会提交。
func (f *QuotaUsageRetryFlow) Execute() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		last = f.state.Next()
		if last == nil {
			return nil
		}
		// 永久拒绝原样返回，避免错误重试后提交。
		if !isTemporaryFailure(last) {
			return last
		}
	}
	return last
}

// isTemporaryFailure 判断 err 是否为可重试的临时占用错误。
// 仅 *RetryFailure 且 Temporary 为 true 才重试，其余错误一律不重试。
func isTemporaryFailure(err error) bool {
	var rf *store.RetryFailure
	if errors.As(err, &rf) {
		return rf.Temporary
	}
	return false
}
