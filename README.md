# Learning Database Internals

A hands-on learning project exploring database internals and distributed systems concepts, based on Martin Kleppmann's distributed systems course.

## Project Structure

```
.
├── martin_kleppmann/          # Course notes and summaries
├── distributed-systems-learning/  # Interactive learning platform
│   ├── apps/
│   │   ├── api/               # Go WebSocket API server
│   │   └── web/               # React frontend (Vite)
│   ├── packages/              # Shared Go packages
│   │   ├── core/              # Core types (nodes, messages, clocks)
│   │   ├── simulation/        # Simulation engine
│   │   ├── visualization/     # Event system for visualization
│   │   ├── failure/           # Failure injection
│   │   ├── network/           # Network transport
│   │   └── protocol/          # WebSocket protocol messages
│   └── projects/              # Distributed systems implementations
│       ├── two-generals/      # Two Generals Problem
│       ├── byzantine/         # Byzantine Fault Tolerance
│       ├── clocks/            # Logical & Vector Clocks
│       ├── broadcast/         # Broadcast Algorithms
│       ├── raft/              # Raft Consensus
│       ├── quorum/            # Quorum Systems
│       ├── state-machine/     # State Machine Replication
│       ├── two-phase-commit/  # 2PC Protocol
│       ├── consistency/       # Consistency Models
│       └── crdt/              # CRDTs
└── employees.sql              # Sample SQL data
```

## Quick Start

### Using Docker (Recommended)

```bash
cd distributed-systems-learning
docker-compose up --build
```

Access the app at http://localhost:3000

### Manual Setup

**Prerequisites:**
- Go 1.23+
- Node.js 22+
- pnpm 8.15+

**Backend:**
```bash
cd distributed-systems-learning
go run ./apps/api/cmd/server/main.go
```

**Frontend:**
```bash
cd distributed-systems-learning
pnpm install
pnpm web:dev
```

## Topics Covered

- Two Generals Problem
- Byzantine Generals Problem
- Physical & Logical Time
- Clock Synchronization
- Causality & Happens-Before
- Broadcast Ordering (FIFO, Causal, Total Order)
- Replication & Quorums
- Consensus (Raft)
- Two-Phase Commit
- Linearizability & Eventual Consistency
- CRDTs

## Resources

- [Martin Kleppmann's Distributed Systems Course](https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/)
- [Designing Data-Intensive Applications](https://dataintensive.net/)
