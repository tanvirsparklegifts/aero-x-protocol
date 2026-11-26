package nexus

import (
	"context"
	"testing"
	"time"

	"aero-x-protocol/pkg/edge"
)

func packet(id string, priority QoSPriority) DataPacket {
	return DataPacket{Payload: edge.TelemetryData{AgentID: id}, Priority: priority, Timestamp: time.Now()}
}

func TestProcessQueueCountsAllPackets(t *testing.T) {
	fabric := NewNexusFabric()
	fabric.Publish(packet("best", BestEffort))
	fabric.Publish(packet("critical", MissionCritical))
	fabric.Publish(packet("reliable", Reliable))
	if got := fabric.ProcessQueue(context.Background()); got != 3 {
		t.Fatalf("ProcessQueue() = %d, want 3", got)
	}
	if got := fabric.ProcessQueue(context.Background()); got != 0 {
		t.Fatalf("second ProcessQueue() = %d, want 0", got)
	}
}

func TestQueueDispatchesPriorityFirstAndPreservesFIFO(t *testing.T) {
	fabric := NewNexusFabricWithCapacity(4)
	fabric.Publish(packet("best-1", BestEffort))
	fabric.Publish(packet("critical-1", MissionCritical))
	fabric.Publish(packet("reliable-1", Reliable))
	fabric.Publish(packet("best-2", BestEffort))

	want := []string{"critical-1", "reliable-1", "best-1", "best-2"}
	for _, expected := range want {
		got, ok := fabric.popNext()
		if !ok || got.Payload.AgentID != expected {
			t.Fatalf("popNext() = %q, %v; want %q", got.Payload.AgentID, ok, expected)
		}
	}
}

func TestPublishRejectsPacketsWhenCapacityReached(t *testing.T) {
	fabric := NewNexusFabricWithCapacity(1)
	if !fabric.Publish(packet("first", BestEffort)) {
		t.Fatal("first packet should fit")
	}
	if fabric.Publish(packet("second", MissionCritical)) {
		t.Fatal("queue should reject packet after reaching capacity")
	}
	if got := fabric.ProcessQueue(context.Background()); got != 1 {
		t.Fatalf("ProcessQueue() = %d, want 1", got)
	}
}

func TestProcessQueueHonorsCancellationAndRetainsPendingPackets(t *testing.T) {
	fabric := NewNexusFabricWithCapacity(3)
	fabric.Publish(packet("one", BestEffort))
	fabric.Publish(packet("two", MissionCritical))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if got := fabric.ProcessQueue(ctx); got != 0 {
		t.Fatalf("cancelled ProcessQueue() = %d, want 0", got)
	}
	if got := fabric.ProcessQueue(context.Background()); got != 2 {
		t.Fatalf("uncancelled ProcessQueue() = %d, want 2", got)
	}
}
