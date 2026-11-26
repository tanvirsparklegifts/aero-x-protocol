# AERO-X Architecture Overview

## Scope

AERO-X is an in-process Go simulation for exploring coordination patterns. It is not flight-control software, an operational aerospace platform, or a validated safety-critical system.

## Component model

```text
cmd/agent (local demo wiring)
       |
       +--> pkg/edge  : telemetry state and snapshots
       |
       +--> pkg/swarm : agent registration, heartbeats,
       |                stale-agent pruning, routing checks
       |
       +--> pkg/nexus : in-memory priority-aware packet queues
```

## Runtime characteristics

- Components run in one process in the current demo.
- Agent state and queue contents are in memory.
- Queue priority ordering is demonstrated locally; no durable delivery guarantee is provided.
- Telemetry snapshots and heartbeat checks are software simulations, not hardware interfaces.
- Local timing output is observational and is not a latency service-level objective.

## Verification

Run the same checks used by the GitHub Actions workflow:

```bash
go test -v ./...
go vet ./...
go build ./...
```

## Future engineering work

1. Add deterministic tests for concurrent telemetry updates and snapshots.
2. Add tests for heartbeat expiry and stale-agent pruning.
3. Add ordering and queue-boundary tests for each packet priority.
4. Define explicit interfaces before introducing network transports or persistence.
5. Document failure handling and resource limits before considering distributed execution.
