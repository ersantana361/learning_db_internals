import { Link } from 'react-router-dom'
import { Database, Search, Lock } from 'lucide-react'

const topics = [
  {
    category: 'Storage Engines',
    icon: Database,
    color: 'var(--accent-storage)',
    items: [
      { id: 'btree', name: 'B-Tree', description: 'Self-balancing tree for sorted data' },
      { id: 'lsm-tree', name: 'LSM Tree', description: 'Write-optimized log-structured storage' },
      { id: 'buffer-pool', name: 'Buffer Pool', description: 'In-memory page cache management' },
      { id: 'page-layout', name: 'Page Layout', description: 'On-disk page structure' },
    ],
  },
  {
    category: 'Query Processing',
    icon: Search,
    color: 'var(--accent-query)',
    items: [
      { id: 'query-parser', name: 'Query Parser', description: 'SQL parsing and AST generation' },
      { id: 'query-optimizer', name: 'Query Optimizer', description: 'Cost-based optimization' },
      { id: 'execution-engine', name: 'Execution Engine', description: 'Query plan execution' },
    ],
  },
  {
    category: 'Transactions',
    icon: Lock,
    color: 'var(--accent-transaction)',
    items: [
      { id: 'mvcc', name: 'MVCC', description: 'Multi-version concurrency control' },
      { id: 'wal', name: 'Write-Ahead Log', description: 'Durability and crash recovery' },
      { id: 'locking', name: 'Locking', description: 'Lock management and deadlocks' },
    ],
  },
]

export default function HomePage() {
  return (
    <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
      <header style={{ textAlign: 'center', marginBottom: '3rem' }}>
        <h1 style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>
          Database Internals
        </h1>
        <p style={{ color: 'var(--text-secondary)', fontSize: '1.1rem' }}>
          Interactive learning platform for understanding how databases work
        </p>
      </header>

      <div style={{ display: 'grid', gap: '2rem' }}>
        {topics.map((section) => (
          <section key={section.category}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '1rem' }}>
              <section.icon size={24} style={{ color: section.color }} />
              <h2 style={{ fontSize: '1.5rem' }}>{section.category}</h2>
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '1rem' }}>
              {section.items.map((topic) => (
                <Link
                  key={topic.id}
                  to={`/topic/${topic.id}`}
                  style={{
                    padding: '1.25rem',
                    background: 'var(--bg-secondary)',
                    borderRadius: '8px',
                    border: '1px solid var(--border-color)',
                    transition: 'border-color 0.2s, transform 0.2s',
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.borderColor = section.color
                    e.currentTarget.style.transform = 'translateY(-2px)'
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.borderColor = 'var(--border-color)'
                    e.currentTarget.style.transform = 'translateY(0)'
                  }}
                >
                  <h3 style={{ marginBottom: '0.5rem' }}>{topic.name}</h3>
                  <p style={{ color: 'var(--text-secondary)', fontSize: '0.9rem' }}>
                    {topic.description}
                  </p>
                </Link>
              ))}
            </div>
          </section>
        ))}
      </div>
    </div>
  )
}
