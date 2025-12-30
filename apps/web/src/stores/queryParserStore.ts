import { create } from 'zustand';
import type { Token, ASTNode, ParsePhase } from '../types';

interface QueryParserStore {
  // Input
  query: string;
  setQuery: (query: string) => void;

  // Tokenization
  tokens: Token[];
  currentTokenIndex: number;
  setTokens: (tokens: Token[]) => void;
  setCurrentTokenIndex: (index: number) => void;
  highlightToken: (index: number) => void;
  clearTokenHighlights: () => void;

  // AST
  astNodes: Record<string, ASTNode>;
  astRoot: string | null;
  setAST: (nodes: Record<string, ASTNode>, rootId: string | null) => void;
  addASTNode: (node: ASTNode) => void;
  highlightASTNode: (nodeId: string) => void;
  clearASTHighlights: () => void;

  // Parse state
  parsePhase: ParsePhase;
  parseError: string | undefined;
  setParsePhase: (phase: ParsePhase) => void;
  setParseError: (error: string | undefined) => void;

  // Reset
  reset: () => void;
}

const initialState = {
  query: '',
  tokens: [],
  currentTokenIndex: -1,
  astNodes: {},
  astRoot: null,
  parsePhase: 'idle' as ParsePhase,
  parseError: undefined,
};

export const useQueryParserStore = create<QueryParserStore>((set) => ({
  ...initialState,

  setQuery: (query) => set({ query }),

  setTokens: (tokens) => set({ tokens }),

  setCurrentTokenIndex: (index) => set({ currentTokenIndex: index }),

  highlightToken: (index) =>
    set((state) => ({
      tokens: state.tokens.map((token, i) => ({
        ...token,
        highlighted: i === index,
      })),
      currentTokenIndex: index,
    })),

  clearTokenHighlights: () =>
    set((state) => ({
      tokens: state.tokens.map((token) => ({ ...token, highlighted: false })),
    })),

  setAST: (nodes, rootId) => set({ astNodes: nodes, astRoot: rootId }),

  addASTNode: (node) =>
    set((state) => ({
      astNodes: { ...state.astNodes, [node.id]: node },
    })),

  highlightASTNode: (nodeId) =>
    set((state) => ({
      astNodes: Object.fromEntries(
        Object.entries(state.astNodes).map(([id, node]) => [
          id,
          { ...node, highlighted: id === nodeId },
        ])
      ),
    })),

  clearASTHighlights: () =>
    set((state) => ({
      astNodes: Object.fromEntries(
        Object.entries(state.astNodes).map(([id, node]) => [
          id,
          { ...node, highlighted: false },
        ])
      ),
    })),

  setParsePhase: (phase) => set({ parsePhase: phase }),

  setParseError: (error) => set({ parseError: error }),

  reset: () => set(initialState),
}));
