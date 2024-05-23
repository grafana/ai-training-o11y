import { create } from 'zustand';

export type ProcessStatus = 'running' | 'finished' | 'crashed' | 'timed out';

export interface RowData {
  process_uuid: string;
  status: ProcessStatus;
  start_time: string;
  end_time?: string;
  [key: string]: string | undefined;
}

interface TrainingAppState {
  // Tab state
  tab: 'table' | 'graphs';
  setTab: (tab: 'table' | 'graphs') => void;

  // Rendered rows state
  renderedRows: RowData[];
  isSelected: boolean[];
  setRenderedRows: (rows: RowData[]) => void;
  setIsSelected: (index: number, value: boolean) => void;

  // Selected rows state
  selectedRows: RowData[];
  indices: Map<string, number | undefined>;
  setSelectedRows: (rows: RowData[]) => void;
  addSelectedRow: (row: RowData) => void;
  removeSelectedRow: (processUuid: string) => void;
}

export const useTrainingAppStore = create<TrainingAppState>()((set) => ({
  // Tab state
  tab: 'table',
  setTab: (tab) => set(() => ({ tab })),

  // Rendered rows state
  renderedRows: [],
  isSelected: [],
  setRenderedRows: (rows) =>
    set(() => {
      const isSelected = new Array(rows.length).fill(false);
      return { renderedRows: rows, isSelected };
    }),
  setIsSelected: (index, value) =>
    set((state) => {
      const newIsSelected = [...state.isSelected];
      newIsSelected[index] = value;
      return { isSelected: newIsSelected };
    }),

  // Selected rows state
  selectedRows: [],
  indices: new Map<string, number | undefined>(),
  setSelectedRows: (rows) =>
    set(() => {
      const newIndices = new Map<string, number | undefined>();
      rows.forEach((row, index) => {
        newIndices.set(row.process_uuid, index);
      });
      return { selectedRows: rows, indices: newIndices };
    }),
  removeSelectedRow: (processUuid) =>
    set((state) => {
      const index = state.indices.get(processUuid);
      if (index !== undefined) {
        const newRows = [...state.selectedRows];
        newRows.splice(index, 1);
        const newIndices = state.indices;
        newIndices.set(processUuid, undefined);
        return { selectedRows: newRows, indices: newIndices };
      }
      return state;
    }),
  addSelectedRow: (row) =>
    set((state) => {
      const newRows = [...state.selectedRows, row];
      const newIndices = state.indices;
      state.indices.set(row.process_uuid, newRows.length - 1);
      return { selectedRows: newRows, indices: newIndices };
    }),
}));
