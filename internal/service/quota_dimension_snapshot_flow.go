package service

import "flowcontrol/internal/store"

type QuotaDimensionSnapshotFlow struct {
	state *store.QuotaDimensionSnapshotState
}

func NewQuotaDimensionSnapshotFlow(state *store.QuotaDimensionSnapshotState) *QuotaDimensionSnapshotFlow {
	return &QuotaDimensionSnapshotFlow{state: state}
}

func (f *QuotaDimensionSnapshotFlow) Publish(key string, values []string) []string {
	f.state.Replace(key, values)
	return f.state.Snapshot(key)
}

func (f *QuotaDimensionSnapshotFlow) Snapshot(key string) []string {
	return f.state.Snapshot(key)
}
