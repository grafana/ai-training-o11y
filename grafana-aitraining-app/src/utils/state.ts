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
  indices: Map<string, number>;
  setRows: (rows: RowData[]) => void;
}

export const useSelectedRowsStore = create<SelectedRowsState>()((set) => ({
  rows: [],
  indices: new Map<string, number>(),
  setRows: (rows) =>
    set(() => {
      const newIndices = new Map<string, number>();
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
        const newIndices = new Map<string, number>();
        newRows.forEach((row, newIndex) => {
          newIndices.set(row.process_uuid, newIndex);
        });
        return { rows: newRows, indices: newIndices };
      }
      return state;
    }),
  }
));
