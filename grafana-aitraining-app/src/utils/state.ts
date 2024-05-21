import { create } from 'zustand'

interface TabState {
  tab: 'table' | 'graphs'
  set: (tab: string) => void
}

export const useTabStore = create<TabState>()((set) => ({
  tab: 'table',
  set: (to) => set((state) => ({ tab: to as 'table' | 'graphs'})),
}))

export type ProcessStatus = 'running' | 'finished' | 'crashed' | 'timed out';

export interface RowData {
  process_uuid: string;
  status: ProcessStatus;
  start_time: string;
  end_time?: string;
  [key: string]: string | undefined;
}

interface RowsState {
  rows: RowData[];
  setRows: (rows: RowData[]) => void;
}

export const useRowsStore = create<RowsState>()((set) => ({
  rows: [],
  setRows: (rows) => set(() => ({ rows })),
}));

interface SelectedRowsState {
  rows: RowData[];
  indices: Map<string, number | undefined>; // Undefined if it was selected previously and is now unselected
  setRows: (rows: RowData[]) => void;
  addRow: (row: RowData) => void;
  removeRow: (processUuid: string) => void;
}

export const useSelectedRowsStore = create<SelectedRowsState>()((set) => ({
  rows: [],
  indices: new Map<string, number | undefined>(),
  setRows: (rows) =>
    set(() => {
      const newIndices = new Map<string, number | undefined>();
      rows.forEach((row, index) => {
        newIndices.set(row.process_uuid, index);
      });
      return { rows, indices: newIndices };
    }),
  removeRow: (processUuid: string) =>
    set((state) => {
      const index = state.indices.get(processUuid);
      if (index !== undefined) {
        const newRows = [...state.rows];
        newRows.splice(index, 1);
        // Instead of the current approach, let us just set that position in the map to undefined
        const newIndices = state.indices;
        newIndices.set(processUuid, undefined);
        return { rows: newRows, indices: newIndices };
      }
      return state;
    }),
  addRow: (row: RowData) =>
    set((state) => {
      const newRows = [...state.rows, row];
      const newIndices = state.indices;
      state.indices.set(row.process_uuid, newRows.length - 1);
      return { rows: newRows, indices: newIndices };
    }),
  }
));
