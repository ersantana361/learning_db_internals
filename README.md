# Database Internals

An interactive learning platform for understanding how databases work under the hood.

## Topics Covered

### Storage Engines
- **B-Tree** - Self-balancing tree data structure for sorted data
- **LSM Tree** - Log-structured merge tree for write-optimized storage
- **Buffer Pool** - In-memory page cache management
- **Page Layout** - On-disk page structure and organization

### Query Processing
- **Query Parser** - SQL parsing and AST generation
- **Query Optimizer** - Cost-based query plan optimization
- **Execution Engine** - Query plan execution and operators

### Transactions
- **MVCC** - Multi-version concurrency control
- **Write-Ahead Log** - Durability and crash recovery
- **Locking** - Lock management and deadlock detection

## Quick Start

### Using Docker

```bash
docker compose up --build
```

Access the app at http://localhost:3001

### Manual Setup

**Prerequisites:**
- Go 1.23+
- Node.js 22+
- pnpm 8.15+

**Backend:**
```bash
go run ./apps/api/cmd/server/main.go
```

**Frontend:**
```bash
pnpm install
pnpm web:dev
```

## Project Structure

```
db_internals/
├── apps/
│   ├── api/                    # Go WebSocket API server
│   └── web/                    # React frontend (Vite)
├── packages/                   # Shared Go packages
│   ├── core/                   # Core types
│   ├── storage/                # Storage engine abstractions
│   ├── query/                  # Query processing
│   ├── transaction/            # Transaction management
│   ├── buffer/                 # Buffer pool
│   └── index/                  # Index structures
├── projects/                   # Interactive implementations
│   ├── btree/                  # B-Tree visualization
│   ├── lsm-tree/               # LSM Tree simulation
│   ├── buffer-pool/            # Buffer pool management
│   ├── page-layout/            # Page structure explorer
│   ├── query-parser/           # SQL parser visualization
│   ├── query-optimizer/        # Query plan optimization
│   ├── execution-engine/       # Execution visualization
│   ├── mvcc/                   # MVCC simulation
│   ├── wal/                    # WAL and recovery
│   └── locking/                # Lock manager simulation
├── docker-compose.yml
├── go.work
└── package.json
```

## Resources

- *Database Internals* by Alex Petrov
- *Designing Data-Intensive Applications* by Martin Kleppmann
- CMU 15-445/645 Database Systems course
