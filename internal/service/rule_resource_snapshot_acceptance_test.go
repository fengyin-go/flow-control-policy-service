package service

import (
	"strings"
	"testing"

	"flowcontrol/internal/store"
)

func TestR001RuleResourceSnapshot(t *testing.T) {
	flow := NewRuleResourceSnapshotFlow(store.NewRuleResourceSnapshotState())
	source := []string{"node-a", "node-b"}
	visible := flow.Publish("blue", source)
	source[0] = "caller-write"
	visible[1] = "response-write"
	got := flow.Snapshot("blue")
	if strings.Join(got, ",") != "node-a,node-b" {
		t.Fatalf("限流规则资源列表 must stay isolated after caller and response mutation: %v", got)
	}
}
