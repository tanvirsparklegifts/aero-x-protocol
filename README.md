# AERO-X Protocol

**Autonomous Aerospace & Multi-Vehicle Orchestration Engine — architecture prototype**

AERO-X is an exploratory Go simulation of three cooperating software components: edge telemetry state, a heartbeat-based agent registry, and a priority-aware in-memory message queue. It is not flight-control software and has not been validated for real aerospace or safety-critical deployment.

## Current implementation

- pkg/edge: mutex-protected telemetry state and snapshot reads.
- pkg/swarm: in-memory agent registration, heartbeat timestamps, stale-agent pruning, and task-routing checks.
- pkg/nexus: in-memory queues that process MissionCritical/Reliable packets before BestEffort packets.
- cmd/agent: a short local simulation wiring the components together.

## Requirements

- Go 1.22+

## Run locally

    go test ./...
    go run ./cmd/agent

## Repository layout

    cmd/agent/             simulation entry point
    pkg/edge/              telemetry state
    pkg/swarm/             agent registry and heartbeat handling
    pkg/nexus/             in-memory priority queue
    proto/telemetry.proto  protocol schema draft
    docs/                  architecture assets

## Scope and limitations

This repository currently demonstrates an in-process simulation. It does not implement real aircraft interfaces, distributed networking, protobuf code generation, durable messaging, real-time guarantees, or flight-certified safety controls. The queue is in-memory and the latency values printed by the demo are observations from a local simulation, not a verified latency SLO.

## Engineering workflow

Pull requests should include focused tests and explain behavior changes. The GitHub Actions workflow runs unit tests and go vet on pushes and pull requests.

See LICENSE for licensing terms.
