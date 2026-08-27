package service

import "flowcontrol/internal/store"

type BreakerProbeSnapshotFlow struct {
	state *store.BreakerProbeSnapshotState
}

func NewBreakerProbeSnapshotFlow(state *store.BreakerProbeSnapshotState) *BreakerProbeSnapshotFlow {
	return &BreakerProbeSnapshotFlow{state: state}
}

func (f *BreakerProbeSnapshotFlow) Publish(key string, values []string) []string {
	f.state.Replace(key, values)
	return f.state.Snapshot(key)
}

func (f *BreakerProbeSnapshotFlow) Snapshot(key string) []string {
	return f.state.Snapshot(key)
}
