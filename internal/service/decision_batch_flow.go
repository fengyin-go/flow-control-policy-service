package service

import "flowcontrol/internal/store"

type DecisionBatchFlow struct{ stream *store.DecisionBatchStream }

func NewDecisionBatchFlow(stream *store.DecisionBatchStream) *DecisionBatchFlow {
	return &DecisionBatchFlow{stream: stream}
}

func (f *DecisionBatchFlow) Collect(items []string, failAt int) ([]string, error) {
	out, errs := f.stream.Start(items, failAt)
	values := make([]string, 0)
	// Start 保证遇错或结束时关闭 out，因此 range 在中断后能及时退出，
	// 已收集的放行结果保留在 values 中，最后读取错误（无错则为 nil）。
	for value := range out {
		values = append(values, value)
	}
	return values, <-errs
}
