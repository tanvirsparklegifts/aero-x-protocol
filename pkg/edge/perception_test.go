package edge

import "testing"

func TestUpdateTelemetryAndSnapshot(t *testing.T) {
	engine := NewEdgeEngine("agent-test")
	got := engine.UpdateTelemetry(23.7, 90.3, 120.0, 42.5)
	if got.AgentID != "agent-test" || got.Latitude != 23.7 || got.Longitude != 90.3 || got.Altitude != 120.0 || got.Velocity != 42.5 {
		t.Fatalf("unexpected telemetry: %+v", got)
	}
	if snapshot := engine.GetCurrentState(); snapshot != got {
		t.Fatalf("snapshot mismatch: got %+v want %+v", snapshot, got)
	}
}
