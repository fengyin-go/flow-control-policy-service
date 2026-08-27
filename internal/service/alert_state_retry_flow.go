package service

import "flowcontrol/internal/store"

// AlertStateRetryFlow 负责将告警状态写入到状态机：遇到临时错误（瞬时繁忙）时重试，
// 遇到永久拒绝时立即中止，不进入下一次写入，避免永久拒绝仍留下提交结果。
type AlertStateRetryFlow struct {
	state *store.AlertStateRetryRetryState
}

func NewAlertStateRetryFlow(state *store.AlertStateRetryRetryState) *AlertStateRetryFlow {
	return &AlertStateRetryFlow{state: state}
}

func (f *AlertStateRetryFlow) Execute() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		last = f.state.Next()
		if last == nil {
			return nil
		}
		// 永久拒绝：重试无意义，且不应让后续写入覆盖本次判定而留下结果。
		if !store.IsTemporary(last) {
			return last
		}
		// 临时错误：进入下一轮重试。
	}
	return last
}
