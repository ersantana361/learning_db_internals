import { create } from 'zustand';
import type { SimulationMode, StepInfo, Highlight } from '../types';

interface TimelineEvent {
  timestamp: number;
  type: string;
  description: string;
}

interface SimulationStore {
  // Connection state
  connected: boolean;
  setConnected: (connected: boolean) => void;

  // Active project
  project: string | null;
  setProject: (project: string | null) => void;

  // Simulation state
  mode: SimulationMode;
  speed: number;
  currentStep: StepInfo | null;
  totalSteps: number;
  setMode: (mode: SimulationMode) => void;
  setSpeed: (speed: number) => void;
  setCurrentStep: (step: StepInfo | null) => void;
  setTotalSteps: (total: number) => void;

  // Steps
  steps: StepInfo[];
  setSteps: (steps: StepInfo[]) => void;

  // Highlights
  highlights: Highlight[];
  setHighlights: (highlights: Highlight[]) => void;
  addHighlight: (highlight: Highlight) => void;
  clearHighlights: () => void;

  // Timeline
  timeline: TimelineEvent[];
  addTimelineEvent: (event: TimelineEvent) => void;
  clearTimeline: () => void;

  // Reset
  reset: () => void;
}

const initialState = {
  connected: false,
  project: null,
  mode: 'idle' as SimulationMode,
  speed: 1.0,
  currentStep: null,
  totalSteps: 0,
  steps: [],
  highlights: [],
  timeline: [],
};

export const useSimulationStore = create<SimulationStore>((set) => ({
  ...initialState,

  setConnected: (connected) => set({ connected }),
  setProject: (project) => set({ project }),
  setMode: (mode) => set({ mode }),
  setSpeed: (speed) => set({ speed }),
  setCurrentStep: (step) => set({ currentStep: step }),
  setTotalSteps: (total) => set({ totalSteps: total }),
  setSteps: (steps) => set({ steps }),

  setHighlights: (highlights) => set({ highlights }),
  addHighlight: (highlight) =>
    set((state) => ({ highlights: [...state.highlights, highlight] })),
  clearHighlights: () => set({ highlights: [] }),

  addTimelineEvent: (event) =>
    set((state) => ({ timeline: [...state.timeline, event] })),
  clearTimeline: () => set({ timeline: [] }),

  reset: () => set(initialState),
}));
