import { create } from 'zustand';

import { PanelData } from '@grafana/data';

export type ProcessStatus = 'empty' | 'running' | 'finished' | 'crashed' | 'timed out';

export interface RowData {
  process_uuid: string;
  status: ProcessStatus;
  start_time: string;
  end_time?: string;
  [key: string]: string | undefined;
}

export interface QueryResultData {
  processData: RowData;
  lokiData: PanelData | undefined; 
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

  // query result state
  queryStatus: ProcessStatus;
  queryData: Record<string, QueryResultData>;
  resetResults: () => void;
  appendResult: (processData: RowData, result: PanelData | undefined) => void;
  setQueryStatus: (status: ProcessStatus) => void;
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
      return { selectedRows: rows };
    }),
  removeSelectedRow: (processUuid) =>
    set((state) => {
      const newRows = state.selectedRows.filter((row) => row.process_uuid !== processUuid);
      return { selectedRows: newRows };
    }),
  addSelectedRow: (row) =>
    set((state) => {
      if (!state.selectedRows.some((selectedRow) => selectedRow.process_uuid === row.process_uuid)) {
        const newRows = [...state.selectedRows, row];
        return { selectedRows: newRows };
      }
      return state;
    }),

    // query result state
  queryData: {},
  queryStatus: 'empty',
  resetResults: () => set(() => ({ queryStatus: 'empty', queryData: {} })),
  appendResult: (processData, result) => set((state) => ({ 
    queryData: { ...state.queryData, [processData.process_uuid]: { processData, lokiData: result } }
  })),
  setQueryStatus: (status: ProcessStatus) => set(() => ({ queryStatus: status })),

}));
