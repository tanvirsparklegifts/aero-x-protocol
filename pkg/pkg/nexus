package nexus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"aero-x-protocol/pkg/edge"
)

// QoSPriority defines transport reliability levels
type QoSPriority int

const (
	BestEffort QoSPriority = iota
	Reliable
	MissionCritical
)

// DataPacket represents an encapsulated message within Nexus Data Fabric
type DataPacket struct {
	Header    string
	Payload   edge.TelemetryData
	Priority  QoSPriority
	Timestamp time.Time
}

// NexusFabric manages QoS-prioritized message queues across swarm nodes
type NexusFabric struct {
	mu           sync.Mutex
	highPriority []DataPacket
	lowPriority  []DataPacket
}

func NewNexusFabric() *NexusFabric {
	return &NexusFabric{
		highPriority: make([]DataPacket, 0),
		lowPriority:  make([]DataPacket, 0),
	}
}

// Publish enqueues data packets based on QoS priority constraints
func (n *NexusFabric) Publish(packet DataPacket) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if packet.Priority == MissionCritical || packet.Priority == Reliable {
		n.highPriority = append(n.highPriority, packet)
		fmt.Printf("[NEXUS DATA FABRIC] 🚨 HIGH PRIORITY Packet Queued from %s | QoS Level: Mission Critical\n", packet.Payload.AgentID)
	} else {
		n.lowPriority = append(n.lowPriority, packet)
		fmt.Printf("[NEXUS DATA FABRIC] ℹ️ Standard Packet Queued from %s | QoS Level: Best Effort\n", packet.Payload.AgentID)
	}
}

// ProcessQueue simulates deterministic, priority-aware message transport
func (n *NexusFabric) ProcessQueue(ctx context.Context) int {
	n.mu.Lock()
	defer n.mu.Unlock()

	processed := 0

	// Process high-priority critical traffic first
	for len(n.highPriority) > 0 {
		pkt := n.highPriority[0]
		n.highPriority = n.highPriority[1:]
		processed++
		fmt.Printf("[NEXUS TRANSPORT] ⚡ Dispatched Critical QoS Telemetry for %s in %v\n", pkt.Payload.AgentID, time.Since(pkt.Timestamp))
	}

	// Process remaining standard traffic
	for len(n.lowPriority) > 0 {
		pkt := n.lowPriority[0]
		n.lowPriority = n.lowPriority[1:]
		processed++
		fmt.Printf("[NEXUS TRANSPORT] 📦 Dispatched Standard Telemetry for %s\n", pkt.Payload.AgentID)
	}

	return processed
}
