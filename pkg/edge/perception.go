package edge

import (
	"sync"
	"time"
)

// TelemetryData represents local perception state vector
type TelemetryData struct {
	AgentID   string    `json:"agent_id"`
	Timestamp time.Time `json:"timestamp"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude"`
	Velocity  float64   `json:"velocity"`
	Health    string    `json:"health"`
}

// EdgeEngine manages local authoritative state
type EdgeEngine struct {
	mu          sync.RWMutex
	currentData TelemetryData
	isOnline    bool
}

func NewEdgeEngine(agentID string) *EdgeEngine {
	return &EdgeEngine{
		currentData: TelemetryData{
			AgentID:   agentID,
			Timestamp: time.Now(),
			Health:    "HEALTHY",
		},
		isOnline: true,
	}
}

// UpdateTelemetry updates vehicle's local perception
func (e *EdgeEngine) UpdateTelemetry(lat, lon, alt, vel float64) TelemetryData {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.currentData.Timestamp = time.Now()
	e.currentData.Latitude = lat
	e.currentData.Longitude = lon
	e.currentData.Altitude = alt
	e.currentData.Velocity = vel

	return e.currentData
}

// GetCurrentState returns local state safely
func (e *EdgeEngine) GetCurrentState() TelemetryData {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.currentData
}
