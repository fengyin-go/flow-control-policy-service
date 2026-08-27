package service

import "flowcontrol/internal/store"

// DecisionLogRetryFlow 负责带重试地写入放行日志。
// 只有临时错误（store.IsTemporary）才重试；遇到永久错误立即停止，不提交。
type DecisionLogRetryFlow struct {
	state *store.DecisionLogRetryRetryState
}

func NewDecisionLogRetryFlow(state *store.DecisionLogRetryRetryState) *DecisionLogRetryFlow {
	return &DecisionLogRetryFlow{state: state}
}

// Execute 执行写入。临时错误最多重试 maxAttempts 次，永久错误立即返回且不提交。
// 重试耗尽仍未成功时，返回最后一次临时错误。
func (f *DecisionLogRetryFlow) Execute() error {
	const maxAttempts = 2
	var last error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		last = f.state.Next()
		switch {
		case last == nil:
			// 写入成功。
			return nil
		case !store.IsTemporary(last):
			// 永久错误：立即停止，不重试、不提交。
			return last
		}
		// 临时错误：继续重试。
	}
	// 重试耗尽：返回最后一次临时错误（Next 在步骤耗尽时也会回退到 lastErr）。
	if last == nil {
		last = f.state.Last()
	}
	return last
}
