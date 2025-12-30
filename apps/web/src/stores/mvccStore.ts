import { create } from 'zustand';
import type { Transaction, Version, Row } from '../types';

interface MVCCStore {
  // Transactions
  transactions: Record<string, Transaction>;
  activeTransaction: string | null;

  // Versions and Rows
  versions: Record<string, Version>;
  rows: Record<string, Row>;

  // Visibility
  currentSnapshot: number;
  visibleVersions: string[];

  // Transaction operations
  setTransactions: (txns: Record<string, Transaction>) => void;
  addTransaction: (txn: Transaction) => void;
  updateTransaction: (txnId: string, update: Partial<Transaction>) => void;
  setActiveTransaction: (txnId: string | null) => void;

  // Version operations
  setVersions: (versions: Record<string, Version>) => void;
  addVersion: (version: Version) => void;

  // Row operations
  setRows: (rows: Record<string, Row>) => void;
  updateRow: (rowId: string, update: Partial<Row>) => void;

  // Visibility
  setSnapshot: (timestamp: number) => void;
  setVisibleVersions: (versions: string[]) => void;
  computeVisibility: () => void;

  // Reset
  reset: () => void;
}

const initialState = {
  transactions: {},
  activeTransaction: null,
  versions: {},
  rows: {},
  currentSnapshot: 0,
  visibleVersions: [],
};

export const useMVCCStore = create<MVCCStore>((set, get) => ({
  ...initialState,

  setTransactions: (transactions) => set({ transactions }),

  addTransaction: (txn) =>
    set((state) => ({
      transactions: { ...state.transactions, [txn.id]: txn },
    })),

  updateTransaction: (txnId, update) =>
    set((state) => ({
      transactions: {
        ...state.transactions,
        [txnId]: { ...state.transactions[txnId], ...update },
      },
    })),

  setActiveTransaction: (txnId) => set({ activeTransaction: txnId }),

  setVersions: (versions) => set({ versions }),

  addVersion: (version) =>
    set((state) => ({
      versions: { ...state.versions, [version.id]: version },
    })),

  setRows: (rows) => set({ rows }),

  updateRow: (rowId, update) =>
    set((state) => ({
      rows: {
        ...state.rows,
        [rowId]: { ...state.rows[rowId], ...update },
      },
    })),

  setSnapshot: (timestamp) => {
    set({ currentSnapshot: timestamp });
    get().computeVisibility();
  },

  setVisibleVersions: (versions) => set({ visibleVersions: versions }),

  computeVisibility: () => {
    const state = get();
    const { versions, transactions, currentSnapshot } = state;

    const visible = Object.values(versions).filter((version) => {
      // Get the creating transaction
      const createdTxn = transactions[version.createdBy];
      if (!createdTxn) return false;

      // Version is visible if:
      // 1. Created by a committed transaction before snapshot
      if (createdTxn.status !== 'committed') return false;
      if (createdTxn.commitTime && createdTxn.commitTime > currentSnapshot)
        return false;

      // 2. Not deleted, or deleted after snapshot
      if (version.deletedBy) {
        const deletedTxn = transactions[version.deletedBy];
        if (
          deletedTxn?.status === 'committed' &&
          deletedTxn.commitTime &&
          deletedTxn.commitTime <= currentSnapshot
        ) {
          return false;
        }
      }

      return true;
    });

    set({ visibleVersions: visible.map((v) => v.id) });
  },

  reset: () => set(initialState),
}));
