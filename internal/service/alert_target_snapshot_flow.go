package service

import "flowcontrol/internal/store"

type AlertTargetSnapshotFlow struct {
	state *store.AlertTargetSnapshotState
}

func NewAlertTargetSnapshotFlow(state *store.AlertTargetSnapshotState) *AlertTargetSnapshotFlow {
	return &AlertTargetSnapshotFlow{state: state}
}

func (f *AlertTargetSnapshotFlow) Publish(key string, values []string) []string {
	f.state.Replace(key, values)
	return f.state.Snapshot(key)
}

func (f *AlertTargetSnapshotFlow) Snapshot(key string) []string {
	return f.state.Snapshot(key)
}
