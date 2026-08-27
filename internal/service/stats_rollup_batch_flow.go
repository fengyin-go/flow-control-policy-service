package service

import "flowcontrol/internal/store"

type StatsRollupBatchFlow struct{ stream *store.StatsRollupBatchStream }

func NewStatsRollupBatchFlow(stream *store.StatsRollupBatchStream) *StatsRollupBatchFlow {
	return &StatsRollupBatchFlow{stream: stream}
}

func (f *StatsRollupBatchFlow) Collect(items []string, failAt int) ([]string, error) {
	out, errs := f.stream.Start(items, failAt)
	values := make([]string, 0)
	for value := range out {
		values = append(values, value)
	}
	return values, <-errs
}
