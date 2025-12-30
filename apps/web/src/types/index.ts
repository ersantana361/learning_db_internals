// Simulation types matching Go protocol package

export type SimulationMode = 'idle' | 'playing' | 'paused' | 'step';

export interface Highlight {
  type: 'node' | 'edge' | 'cell' | 'row' | 'token';
  id: string;
  color: string;
  animation?: 'pulse' | 'flash' | 'fade' | 'fadeIn' | 'shake' | 'none';
  duration?: number;
}

export interface StepInfo {
  index: number;
  total: number;
  title: string;
  description: string;
  highlights: Highlight[];
  data?: Record<string, unknown>;
}

export interface SimulationState {
  project: string;
  mode: SimulationMode;
  speed: number;
  currentStep: StepInfo | null;
  totalSteps: number;
  data: Record<string, unknown>;
}

export interface SimulationConfig {
  project: string;
  scenario?: string;
  stepMode?: boolean;
  speed?: number;
  parameters?: Record<string, unknown>;
}

// Message types
export type MessageType =
  | 'start_simulation'
  | 'pause_simulation'
  | 'resume_simulation'
  | 'stop_simulation'
  | 'step_forward'
  | 'step_backward'
  | 'set_speed'
  | 'reset'
  | 'execute_operation'
  | 'select_scenario'
  | 'get_state'
  | 'simulation_state'
  | 'step_update'
  | 'node_highlight'
  | 'node_update'
  | 'error';

export interface Message<T = unknown> {
  type: MessageType;
  payload?: T;
}

// B-Tree types
export interface BTreeNode {
  id: string;
  keys: number[];
  children: string[];
  isLeaf: boolean;
  parent: string | null;
  highlighted?: boolean;
  highlightColor?: string;
}

export interface BTreeState {
  nodes: Record<string, BTreeNode>;
  rootId: string | null;
  order: number;
  currentPath: string[];
  stats: {
    height: number;
    nodeCount: number;
    keyCount: number;
    comparisons: number;
  };
}

// MVCC types
export type TransactionStatus = 'active' | 'committed' | 'aborted';

export interface Transaction {
  id: string;
  startTime: number;
  commitTime?: number;
  status: TransactionStatus;
  readSet: string[];
  writeSet: string[];
}

export interface Version {
  id: string;
  rowId: string;
  data: Record<string, unknown>;
  createdBy: string;
  createdAt: number;
  deletedBy?: string;
  deletedAt?: number;
  prev?: string;
}

export interface Row {
  id: string;
  currentVersion: string;
  versionChain: string[];
}

export interface MVCCState {
  transactions: Record<string, Transaction>;
  versions: Record<string, Version>;
  rows: Record<string, Row>;
  activeTransaction: string | null;
  currentSnapshot: number;
  visibleVersions: string[];
}

// Query Parser types
export interface Token {
  type: string;
  value: string;
  position: { start: number; end: number };
  highlighted?: boolean;
}

export interface ASTNode {
  id: string;
  type: string;
  value?: string;
  children: string[];
  parent?: string;
  highlighted?: boolean;
}

export type ParsePhase = 'idle' | 'tokenizing' | 'parsing' | 'complete' | 'error';

export interface QueryParserState {
  query: string;
  tokens: Token[];
  currentTokenIndex: number;
  astNodes: Record<string, ASTNode>;
  astRoot: string | null;
  parsePhase: ParsePhase;
  parseError?: string;
}
