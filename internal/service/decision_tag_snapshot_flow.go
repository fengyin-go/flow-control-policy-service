package service

import "flowcontrol/internal/store"

type DecisionTagSnapshotFlow struct {
	state *store.DecisionTagSnapshotState
}

func NewDecisionTagSnapshotFlow(state *store.DecisionTagSnapshotState) *DecisionTagSnapshotFlow {
	return &DecisionTagSnapshotFlow{state: state}
}

func (f *DecisionTagSnapshotFlow) Publish(key string, values []string) []string {
	f.state.Replace(key, values)
	return f.state.Snapshot(key)
}

func (f *DecisionTagSnapshotFlow) Snapshot(key string) []string {
	return f.state.Snapshot(key)
}
