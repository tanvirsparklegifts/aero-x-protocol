package nexus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"aero-x-protocol/pkg/edge"
)

type QoSPriority int

const (
	BestEffort QoSPriority = iota
	Reliable
	MissionCritical
)

const defaultQueueCapacity = 1024

type DataPacket struct {
	Header    string
	Payload   edge.TelemetryData
	Priority  QoSPriority
	Timestamp time.Time
}

// NexusFabric is an in-memory, bounded, priority-aware queue for simulation.
type NexusFabric struct {
	mu           sync.Mutex
	highPriority []DataPacket
	lowPriority  []DataPacket
	maxQueueSize int
}

func NewNexusFabric() *NexusFabric {
	return NewNexusFabricWithCapacity(defaultQueueCapacity)
}

// NewNexusFabricWithCapacity creates a queue with a fixed total capacity.
func NewNexusFabricWithCapacity(capacity int) *NexusFabric {
	if capacity < 1 {
		capacity = defaultQueueCapacity
	}
	return &NexusFabric{
		highPriority: make([]DataPacket, 0),
		lowPriority:  make([]DataPacket, 0),
		maxQueueSize: capacity,
	}
}

// Publish returns false when the bounded queue is full.
func (n *NexusFabric) Publish(packet DataPacket) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	if len(n.highPriority)+len(n.lowPriority) >= n.maxQueueSize {
		return false
	}
	if packet.Priority == MissionCritical || packet.Priority == Reliable {
		n.highPriority = append(n.highPriority, packet)
		fmt.Printf("[NEXUS DATA FABRIC] HIGH PRIORITY packet queued from %s\n", packet.Payload.AgentID)
	} else {
		n.lowPriority = append(n.lowPriority, packet)
		fmt.Printf("[NEXUS DATA FABRIC] Best-effort packet queued from %s\n", packet.Payload.AgentID)
	}
	return true
}

// popNext removes one packet, always preferring the high-priority FIFO.
func (n *NexusFabric) popNext() (DataPacket, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if len(n.highPriority) > 0 {
		packet := n.highPriority[0]
		n.highPriority[0] = DataPacket{}
		n.highPriority = n.highPriority[1:]
		return packet, true
	}
	if len(n.lowPriority) > 0 {
		packet := n.lowPriority[0]
		n.lowPriority[0] = DataPacket{}
		n.lowPriority = n.lowPriority[1:]
		return packet, true
	}
	return DataPacket{}, false
}

// ProcessQueue drains packets until empty or ctx is cancelled. Unprocessed
// packets remain queued when cancellation occurs.
func (n *NexusFabric) ProcessQueue(ctx context.Context) int {
	processed := 0
	for {
		if err := ctx.Err(); err != nil {
			return processed
		}
		packet, ok := n.popNext()
		if !ok {
			return processed
		}
		processed++
		if packet.Priority == MissionCritical || packet.Priority == Reliable {
			fmt.Printf("[NEXUS TRANSPORT] Dispatched priority telemetry for %s in %v\n", packet.Payload.AgentID, time.Since(packet.Timestamp))
		} else {
			fmt.Printf("[NEXUS TRANSPORT] Dispatched best-effort telemetry for %s\n", packet.Payload.AgentID)
		}
	}
}
