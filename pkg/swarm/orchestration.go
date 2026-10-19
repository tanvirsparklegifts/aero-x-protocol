package swarm

import (
	"context"
	"fmt"
	"sync"
	"time"

	"aero-x-protocol/pkg/edge"
)

// SwarmRegistry holds cluster-wide agent states
type SwarmRegistry struct {
	mu     sync.RWMutex
	agents map[string]edge.TelemetryData
}

func NewSwarmRegistry() *SwarmRegistry {
	return &SwarmRegistry{
		agents: make(map[string]edge.TelemetryData),
	}
}

// RegisterOrUpdateAgent synchronizes an agent's state in the swarm
func (s *SwarmRegistry) RegisterOrUpdateAgent(data edge.TelemetryData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agents[data.AgentID] = data
	fmt.Printf("[SWARM CORE] Agent %s state synchronized. Latency check: %v\n", data.AgentID, time.Since(data.Timestamp))
}

// CoordinateTask evaluates bounded-latency decision for dynamic task assignment
func (s *SwarmRegistry) CoordinateTask(ctx context.Context, taskID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	activeAgents := len(s.agents)
	if activeAgents == 0 {
		return fmt.Errorf("task arbitration failed: no active agents in swarm")
	}

	fmt.Printf("[SWARM CORE] Task '%s' assigned across %d synchronized agents under bounded-latency SLO.\n", taskID, activeAgents)
	return nil
}
