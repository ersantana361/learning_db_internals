interface Project {
  id: string;
  name: string;
  description: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
}

const projects: Project[] = [
  {
    id: 'two-generals',
    name: 'Two Generals Problem',
    description: 'Explore message acknowledgment and unreliable channels',
    difficulty: 'beginner',
  },
  {
    id: 'clocks',
    name: 'Logical & Physical Clocks',
    description: 'Understand Lamport and Vector clocks',
    difficulty: 'beginner',
  },
  {
    id: 'byzantine',
    name: 'Byzantine Generals',
    description: 'Handle malicious actors with 3f+1 fault tolerance',
    difficulty: 'intermediate',
  },
  {
    id: 'broadcast',
    name: 'Broadcast Protocols',
    description: 'FIFO, Causal, and Total Order broadcast',
    difficulty: 'intermediate',
  },
  {
    id: 'quorum',
    name: 'Quorum Systems',
    description: 'Read/Write quorums with w+r>n guarantee',
    difficulty: 'intermediate',
  },
  {
    id: 'state-machine',
    name: 'State Machine Replication',
    description: 'Replicated logs with deterministic state',
    difficulty: 'intermediate',
  },
  {
    id: 'raft',
    name: 'Raft Consensus',
    description: 'Leader election and log replication',
    difficulty: 'advanced',
  },
  {
    id: 'two-phase-commit',
    name: 'Two-Phase Commit',
    description: 'Distributed transactions with atomic commits',
    difficulty: 'advanced',
  },
  {
    id: 'consistency',
    name: 'Consistency Models',
    description: 'Linearizability vs Eventual Consistency',
    difficulty: 'advanced',
  },
  {
    id: 'crdt',
    name: 'CRDTs',
    description: 'Conflict-free collaborative editing',
    difficulty: 'advanced',
  },
];

interface ProjectSelectorProps {
  onSelect: (projectId: string) => void;
  currentProject: string | null;
}

export function ProjectSelector({ onSelect, currentProject }: ProjectSelectorProps) {
  const difficultyColors = {
    beginner: 'bg-green-100 text-green-800',
    intermediate: 'bg-yellow-100 text-yellow-800',
    advanced: 'bg-red-100 text-red-800',
  };

  return (
    <div className="project-selector">
      <h2>Select a Project</h2>
      <div className="project-grid">
        {projects.map((project) => (
          <div
            key={project.id}
            className={`project-card ${currentProject === project.id ? 'selected' : ''}`}
            onClick={() => onSelect(project.id)}
          >
            <div className="project-header">
              <h3>{project.name}</h3>
              <span className={`difficulty-badge ${difficultyColors[project.difficulty]}`}>
                {project.difficulty}
              </span>
            </div>
            <p>{project.description}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
