package service

import "flowcontrol/internal/store"

type RuleResourceSnapshotFlow struct {
	state *store.RuleResourceSnapshotState
}

func NewRuleResourceSnapshotFlow(state *store.RuleResourceSnapshotState) *RuleResourceSnapshotFlow {
	return &RuleResourceSnapshotFlow{state: state}
}

func (f *RuleResourceSnapshotFlow) Publish(key string, values []string) []string {
	f.state.Replace(key, values)
	return f.state.Snapshot(key)
}

func (f *RuleResourceSnapshotFlow) Snapshot(key string) []string {
	return f.state.Snapshot(key)
}
