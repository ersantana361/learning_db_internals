import { useParams, Link } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'

const topicDetails: Record<string, { name: string; description: string; content: string }> = {
  btree: {
    name: 'B-Tree',
    description: 'Self-balancing tree data structure',
    content: `
B-Trees are the backbone of most database index structures. They maintain sorted data
and allow searches, insertions, and deletions in logarithmic time.

Key concepts:
- Balanced tree structure with a configurable branching factor
- All leaves at the same depth
- Nodes contain multiple keys and child pointers
- Efficient for disk-based storage due to minimizing I/O operations
    `,
  },
  'lsm-tree': {
    name: 'LSM Tree',
    description: 'Log-structured merge tree',
    content: `
LSM Trees optimize for write-heavy workloads by buffering writes in memory
and periodically flushing to disk in sorted runs.

Key concepts:
- Memtable for in-memory writes
- SSTables (Sorted String Tables) on disk
- Compaction to merge and sort data
- Write amplification vs read amplification tradeoffs
    `,
  },
  'buffer-pool': {
    name: 'Buffer Pool',
    description: 'In-memory page cache',
    content: `
The buffer pool is a cache of database pages in memory, reducing disk I/O
by keeping frequently accessed pages readily available.

Key concepts:
- Page replacement policies (LRU, Clock, etc.)
- Dirty page tracking
- Pin counting for concurrent access
- Prefetching strategies
    `,
  },
  'page-layout': {
    name: 'Page Layout',
    description: 'On-disk page structure',
    content: `
Database pages are fixed-size blocks that store tuples and metadata.
Understanding page layout is crucial for storage efficiency.

Key concepts:
- Slotted page structure
- Tuple headers and data
- Free space management
- Page headers and checksums
    `,
  },
  'query-parser': {
    name: 'Query Parser',
    description: 'SQL parsing',
    content: `
The query parser transforms SQL text into an Abstract Syntax Tree (AST)
that can be analyzed and optimized.

Key concepts:
- Lexical analysis (tokenization)
- Syntax analysis (parsing)
- AST representation
- Semantic analysis
    `,
  },
  'query-optimizer': {
    name: 'Query Optimizer',
    description: 'Cost-based optimization',
    content: `
The query optimizer chooses the most efficient execution plan
based on statistics and cost models.

Key concepts:
- Cost estimation
- Join ordering
- Index selection
- Plan enumeration
    `,
  },
  'execution-engine': {
    name: 'Execution Engine',
    description: 'Query execution',
    content: `
The execution engine runs the optimized query plan using
operators like scans, joins, and aggregations.

Key concepts:
- Volcano/Iterator model
- Vectorized execution
- Pipeline breakers
- Memory management
    `,
  },
  mvcc: {
    name: 'MVCC',
    description: 'Multi-version concurrency control',
    content: `
MVCC allows multiple transactions to read and write data concurrently
without blocking each other by maintaining multiple versions.

Key concepts:
- Transaction IDs and visibility
- Snapshot isolation
- Version chains
- Garbage collection
    `,
  },
  wal: {
    name: 'Write-Ahead Log',
    description: 'Durability and recovery',
    content: `
The WAL ensures durability by logging changes before they're applied,
enabling crash recovery.

Key concepts:
- Log sequence numbers (LSN)
- ARIES recovery algorithm
- Checkpointing
- Log truncation
    `,
  },
  locking: {
    name: 'Locking',
    description: 'Concurrency control',
    content: `
Lock-based concurrency control prevents conflicts between
concurrent transactions.

Key concepts:
- Lock modes (shared, exclusive)
- Two-phase locking (2PL)
- Deadlock detection
- Lock granularity
    `,
  },
}

export default function TopicPage() {
  const { id } = useParams<{ id: string }>()
  const topic = id ? topicDetails[id] : null

  if (!topic) {
    return (
      <div style={{ padding: '2rem', textAlign: 'center' }}>
        <h1>Topic not found</h1>
        <Link to="/">Go back home</Link>
      </div>
    )
  }

  return (
    <div style={{ padding: '2rem', maxWidth: '900px', margin: '0 auto' }}>
      <Link
        to="/"
        style={{
          display: 'inline-flex',
          alignItems: 'center',
          gap: '0.5rem',
          color: 'var(--text-secondary)',
          marginBottom: '2rem',
        }}
      >
        <ArrowLeft size={20} />
        Back to topics
      </Link>

      <h1 style={{ fontSize: '2rem', marginBottom: '0.5rem' }}>{topic.name}</h1>
      <p style={{ color: 'var(--text-secondary)', marginBottom: '2rem' }}>
        {topic.description}
      </p>

      <div
        style={{
          background: 'var(--bg-secondary)',
          padding: '1.5rem',
          borderRadius: '8px',
          whiteSpace: 'pre-wrap',
          lineHeight: '1.8',
        }}
      >
        {topic.content}
      </div>

      <div
        style={{
          marginTop: '2rem',
          padding: '1.5rem',
          background: 'var(--bg-tertiary)',
          borderRadius: '8px',
          textAlign: 'center',
        }}
      >
        <p style={{ color: 'var(--text-secondary)' }}>
          Interactive visualization coming soon...
        </p>
      </div>
    </div>
  )
}
