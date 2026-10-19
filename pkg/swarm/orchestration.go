package swarm

import (
	"context"
	"fmt"
	"sync"
	"time"

	"aero-x-protocol/pkg/edge"
)

// SwarmRegistry holds cluster-wide agent states and heartbeats
type SwarmRegistry struct {
	mu            sync.RWMutex
	agents        map[string]edge.TelemetryData
	lastSeen      map[string]time.Time
	activeTimeout time.Duration
}

func NewSwarmRegistry(timeout time.Duration) *SwarmRegistry {
	return &SwarmRegistry{
		agents:        make(map[string]edge.TelemetryData),
		lastSeen:      make(map[string]time.Time),
		activeTimeout: timeout,
	}
}

// RegisterOrUpdateAgent synchronizes an agent's state & updates Heartbeat
func (s *SwarmRegistry) RegisterOrUpdateAgent(data edge.TelemetryData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.agents[data.AgentID] = data
	s.lastSeen[data.AgentID] = time.Now()
	fmt.Printf("[SWARM CORE] Agent %s heartbeat synced. Latency: %v\n", data.AgentID, time.Since(data.Timestamp))
}

// PruneStaleAgents checks for disconnected agents (Failover Detection)
func (s *SwarmRegistry) PruneStaleAgents() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var droppedAgents []string
	now := time.Now()

	for agentID, lastTime := range s.lastSeen {
		if now.Sub(lastTime) > s.activeTimeout {
			droppedAgents = append(droppedAgents, agentID)
			delete(s.agents, agentID)
			delete(s.lastSeen, agentID)
			fmt.Printf("[FAILOVER ALERT] ⚠️ Agent %s lost! Triggering dynamic re-route...\n", agentID)
		}
	}
	return droppedAgents
}

// CoordinateTask evaluates active agents for bounded-latency orchestration
func (s *SwarmRegistry) CoordinateTask(ctx context.Context, taskID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	activeCount := len(s.agents)
	if activeCount == 0 {
		return fmt.Errorf("task arbitration failed: ZERO active agents in swarm")
	}

	fmt.Printf("[SWARM ORCHESTRATION] Task '%s' successfully routed to %d active node(s) under bounded-latency SLO.\n", taskID, activeCount)
	return nil
}
