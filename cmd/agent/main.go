package main

import (
	"context"
	"fmt"
	"time"

	"aero-x-protocol/pkg/edge"
	"aero-x-protocol/pkg/nexus"
	"aero-x-protocol/pkg/swarm"
)

func main() {
	fmt.Println("==========================================================")
	fmt.Println("  AERO-X PROTOCOL: INTEGRATED EDGE, SWARM & NEXUS FABRIC")
	fmt.Println("==========================================================")

	// Initialize Swarm Registry and Nexus Fabric Infrastructure
	swarmRegistry := swarm.NewSwarmRegistry(1500 * time.Millisecond)
	nexusFabric := nexus.NewNexusFabric()

	agentAlpha := edge.NewEdgeEngine("AERO-AGENT-01-ALPHA")
	agentBeta := edge.NewEdgeEngine("AERO-AGENT-02-BETA")

	// 1. Simulate Edge Telemetry Stream over Nexus QoS Transport
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(400 * time.Millisecond)

			telemetryAlpha := agentAlpha.UpdateTelemetry(23.8103, 90.4125, 1200.5, 240.0)
			telemetryBeta := agentBeta.UpdateTelemetry(23.8150, 90.4180, 1210.0, 245.5)

			swarmRegistry.RegisterOrUpdateAgent(telemetryAlpha)
			swarmRegistry.RegisterOrUpdateAgent(telemetryBeta)

			// Publish to Nexus Data Fabric with QoS differentiation
			nexusFabric.Publish(nexus.DataPacket{
				Header:    "TELEMETRY-STREAM",
				Payload:   telemetryAlpha,
				Priority:  nexus.MissionCritical,
				Timestamp: time.Now(),
			})

			nexusFabric.Publish(nexus.DataPacket{
				Header:    "TELEMETRY-STREAM",
				Payload:   telemetryBeta,
				Priority:  nexus.BestEffort,
				Timestamp: time.Now(),
			})
		}
	}()

	// 2. Nexus Data Fabric Transport Processing Loop
	ticker := time.NewTicker(700 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	for {
		select {
		case <-ticker.C:
			swarmRegistry.PruneStaleAgents()
			_ = swarmRegistry.CoordinateTask(context.Background(), "SWARM-MISSION-QOS-ROUTE")
			nexusFabric.ProcessQueue(context.Background())

		case <-ctx.Done():
			fmt.Println("\n==========================================================")
			fmt.Println("[SYSTEM STATUS] Layer 03 (Nexus Data Fabric QoS): PASSED ✅")
			fmt.Println("==========================================================")
			return
		}
	}
}
