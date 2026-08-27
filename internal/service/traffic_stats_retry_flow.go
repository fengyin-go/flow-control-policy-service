package service

import "flowcontrol/internal/store"

type TrafficStatsRetryFlow struct {
	state *store.TrafficStatsRetryRetryState
}

func NewTrafficStatsRetryFlow(state *store.TrafficStatsRetryRetryState) *TrafficStatsRetryFlow {
	return &TrafficStatsRetryFlow{state: state}
}

func (f *TrafficStatsRetryFlow) Execute() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		last = f.state.Next()
		if last == nil {
			return nil
		}
	}
	return last
}
