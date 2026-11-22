package nexus

import (
	"context"
	"testing"
	"time"

	"aero-x-protocol/pkg/edge"
)

func TestProcessQueueCountsAllPackets(t *testing.T) {
	fabric := NewNexusFabric()
	payload := edge.TelemetryData{AgentID: "agent-test"}
	fabric.Publish(DataPacket{Payload: payload, Priority: BestEffort, Timestamp: time.Now()})
	fabric.Publish(DataPacket{Payload: payload, Priority: MissionCritical, Timestamp: time.Now()})
	fabric.Publish(DataPacket{Payload: payload, Priority: Reliable, Timestamp: time.Now()})
	if got := fabric.ProcessQueue(context.Background()); got != 3 {
		t.Fatalf("ProcessQueue() = %d, want 3", got)
	}
	if got := fabric.ProcessQueue(context.Background()); got != 0 {
		t.Fatalf("second ProcessQueue() = %d, want 0", got)
	}
}
