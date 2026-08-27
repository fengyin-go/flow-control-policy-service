package service

import "flowcontrol/internal/store"

type TrafficStatsRetryFlow struct {
	state *store.TrafficStatsRetryRetryState
}

func NewTrafficStatsRetryFlow(state *store.TrafficStatsRetryRetryState) *TrafficStatsRetryFlow {
	return &TrafficStatsRetryFlow{state: state}
}

// Execute 执行流量统计写入：临时故障按上限重试，永久拒绝立即停止。
func (f *TrafficStatsRetryFlow) Execute() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		last = f.state.Next()
		if last == nil {
			return nil
		}
		// 永久拒绝不可重试：立即停止并返回原始拒绝原因，
		// 避免被当成临时故障重复执行，或在后续轮次被错误提交。
		if store.IsPermanentFailure(last) {
			return last
		}
		// 临时故障：继续重试。
	}
	return last
}
