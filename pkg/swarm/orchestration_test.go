package swarm

import (
	"context"
	"testing"
	"time"

	"aero-x-protocol/pkg/edge"
)

func TestRegisterAndPruneStaleAgent(t *testing.T) {
	registry := NewSwarmRegistry(5 * time.Millisecond)
	registry.RegisterOrUpdateAgent(edge.TelemetryData{AgentID: "agent-test", Timestamp: time.Now()})
	if err := registry.CoordinateTask(context.Background(), "test-task"); err != nil {
		t.Fatalf("expected active agent to accept task: %v", err)
	}
	time.Sleep(10 * time.Millisecond)
	dropped := registry.PruneStaleAgents()
	if len(dropped) != 1 || dropped[0] != "agent-test" {
		t.Fatalf("unexpected pruned agents: %v", dropped)
	}
	if err := registry.CoordinateTask(context.Background(), "test-task"); err == nil {
		t.Fatal("expected task routing to fail with no active agents")
	}
}
