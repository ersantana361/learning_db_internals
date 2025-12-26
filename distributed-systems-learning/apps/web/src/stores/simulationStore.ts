import { create } from 'zustand';

export interface NodeState {
  id: string;
  status: 'running' | 'crashed' | 'partitioned' | 'byzantine';
  role?: string;
  term?: number;
  votedFor?: string;
  log?: LogEntry[];
  commitIndex?: number;
  clock?: Record<string, number>;
  customState?: Record<string, any>;
}

export interface LogEntry {
  index: number;
  term: number;
  command: any;
}

export interface TimelineEvent {
  type: string;
  time: number;
  data: Record<string, any>;
}

export interface SimulationState {
  mode: 'realtime' | 'step' | 'paused';
  speed: number;
  virtualTime: number;
  running: boolean;
}

export interface Partition {
  from: string;
  to: string;
}

interface SimulationStore {
  // Connection state
  connected: boolean;
  setConnected: (connected: boolean) => void;

  // Simulation state
  simulation: SimulationState;
  setSimulationState: (state: Partial<SimulationState>) => void;

  // Node state
  nodes: Record<string, NodeState>;
  setNodes: (nodes: Record<string, NodeState>) => void;
  updateNode: (nodeId: string, update: Partial<NodeState>) => void;

  // Timeline events
  timeline: TimelineEvent[];
  addTimelineEvent: (event: TimelineEvent) => void;
  clearTimeline: () => void;

  // Partitions
  partitions: Partition[];
  addPartition: (partition: Partition) => void;
  removePartition: (from: string, to: string) => void;
  clearPartitions: () => void;

  // Current project
  currentProject: string | null;
  setCurrentProject: (project: string | null) => void;

  // UI state
  selectedNode: string | null;
  setSelectedNode: (nodeId: string | null) => void;

  // Reset
  reset: () => void;
}

const initialSimulationState: SimulationState = {
  mode: 'paused',
  speed: 1.0,
  virtualTime: 0,
  running: false,
};

export const useSimulationStore = create<SimulationStore>((set) => ({
  // Connection
  connected: false,
  setConnected: (connected) => set({ connected }),

  // Simulation
  simulation: initialSimulationState,
  setSimulationState: (state) =>
    set((s) => ({ simulation: { ...s.simulation, ...state } })),

  // Nodes
  nodes: {},
  setNodes: (nodes) => set({ nodes }),
  updateNode: (nodeId, update) =>
    set((s) => ({
      nodes: {
        ...s.nodes,
        [nodeId]: { ...s.nodes[nodeId], ...update },
      },
    })),

  // Timeline
  timeline: [],
  addTimelineEvent: (event) =>
    set((s) => ({
      timeline: [...s.timeline.slice(-999), event], // Keep last 1000 events
    })),
  clearTimeline: () => set({ timeline: [] }),

  // Partitions
  partitions: [],
  addPartition: (partition) =>
    set((s) => ({
      partitions: [...s.partitions, partition],
    })),
  removePartition: (from, to) =>
    set((s) => ({
      partitions: s.partitions.filter(
        (p) => !(p.from === from && p.to === to)
      ),
    })),
  clearPartitions: () => set({ partitions: [] }),

  // Current project
  currentProject: null,
  setCurrentProject: (project) => set({ currentProject: project }),

  // UI
  selectedNode: null,
  setSelectedNode: (nodeId) => set({ selectedNode: nodeId }),

  // Reset
  reset: () =>
    set({
      simulation: initialSimulationState,
      nodes: {},
      timeline: [],
      partitions: [],
      currentProject: null,
      selectedNode: null,
    }),
}));
