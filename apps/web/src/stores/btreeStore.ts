import { create } from 'zustand';
import type { BTreeNode } from '../types';

// Export node data type for visualization components
export type BTreeNodeData = BTreeNode;

interface BTreeStats {
  height: number;
  nodeCount: number;
  keyCount: number;
  comparisons: number;
}

interface BTreeStore {
  // Tree structure
  nodes: Record<string, BTreeNode>;
  rootId: string | null;
  order: number;

  // Operations
  setNodes: (nodes: Record<string, BTreeNode>) => void;
  updateNode: (nodeId: string, update: Partial<BTreeNode>) => void;
  setRoot: (rootId: string | null) => void;
  setOrder: (order: number) => void;

  // Traversal state
  currentPath: string[];
  setCurrentPath: (path: string[]) => void;

  // Statistics
  stats: BTreeStats;
  setStats: (stats: Partial<BTreeStats>) => void;

  // Highlight management
  highlightNode: (nodeId: string, color?: string) => void;
  clearNodeHighlights: () => void;

  // Reset
  reset: () => void;
}

const initialStats: BTreeStats = {
  height: 0,
  nodeCount: 0,
  keyCount: 0,
  comparisons: 0,
};

const initialState = {
  nodes: {},
  rootId: null,
  order: 3,
  currentPath: [],
  stats: initialStats,
};

export const useBTreeStore = create<BTreeStore>((set) => ({
  ...initialState,

  setNodes: (nodes) => set({ nodes }),

  updateNode: (nodeId, update) =>
    set((state) => ({
      nodes: {
        ...state.nodes,
        [nodeId]: { ...state.nodes[nodeId], ...update },
      },
    })),

  setRoot: (rootId) => set({ rootId }),
  setOrder: (order) => set({ order }),
  setCurrentPath: (path) => set({ currentPath: path }),

  setStats: (stats) =>
    set((state) => ({ stats: { ...state.stats, ...stats } })),

  highlightNode: (nodeId, color = '#fbbf24') =>
    set((state) => ({
      nodes: {
        ...state.nodes,
        [nodeId]: {
          ...state.nodes[nodeId],
          highlighted: true,
          highlightColor: color,
        },
      },
    })),

  clearNodeHighlights: () =>
    set((state) => ({
      nodes: Object.fromEntries(
        Object.entries(state.nodes).map(([id, node]) => [
          id,
          { ...node, highlighted: false, highlightColor: undefined },
        ])
      ),
    })),

  reset: () => set(initialState),
}));
