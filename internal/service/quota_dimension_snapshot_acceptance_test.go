package service

import (
	"strings"
	"testing"

	"flowcontrol/internal/store"
)

func TestR002QuotaDimensionSnapshot(t *testing.T) {
	flow := NewQuotaDimensionSnapshotFlow(store.NewQuotaDimensionSnapshotState())
	source := []string{"node-a", "node-b"}
	visible := flow.Publish("blue", source)
	source[0] = "caller-write"
	visible[1] = "response-write"
	got := flow.Snapshot("blue")
	if strings.Join(got, ",") != "node-a,node-b" {
		t.Fatalf("配额维度列表 must stay isolated after caller and response mutation: %v", got)
	}
}
