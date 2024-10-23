import { create } from 'zustand';
import { PanelData } from '@grafana/data';

export type ConfigStatus = 'loading' | 'success' | 'error';
export type ProcessStatus = 'empty' | 'running' | 'finished' | 'crashed' | 'timed out';
export type QueryStatus =
  | 'idle'
  | 'loading'
  | 'success'
  | 'error'
  | 'unauthorized'
  | 'notFound'
  | 'serverError';

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

  // Process table query state
  processesQueryStatus: QueryStatus;
  setProcessesQueryStatus: (status: QueryStatus) => void;

  // Rendered rows state
  renderedRows: RowData[];
  isSelected: boolean[];
  setRenderedRows: (rows: RowData[]) => void;
  setIsSelected: (index: number, value: boolean) => void;

  // Selected rows state
  selectedRows: RowData[];
  selectedRowsMap: Record<string, number>;
  setSelectedRows: (rows: RowData[]) => void;
  addSelectedRow: (row: RowData) => void;
  removeSelectedRow: (processUuid: string) => void;
}

export const useTrainingAppStore = create<TrainingAppState>()((set) => ({
  // Tab state
  tab: 'table',
  setTab: (tab) => set(() => ({ tab })),

  processesQueryStatus: 'idle',
  setProcessesQueryStatus: (status: QueryStatus) => set(() => ({ processesQueryStatus: status })),

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
  selectedRowsMap: {},
  setSelectedRows: (rows) =>
    set(() => {
      const selectedRowsMap: Record<string, number> = {};
      rows.forEach((row, index) => {
        selectedRowsMap[row.process_uuid] = index;
      });
      return { selectedRows: rows, selectedRowsMap };
    }),
  removeSelectedRow: (processUuid) =>
    set((state) => {
      const index = state.selectedRowsMap[processUuid];
      if (index !== undefined) {
        const newSelectedRows = [...state.selectedRows];
        newSelectedRows.splice(index, 1);
        const newSelectedRowsMap: Record<string, number> = {};
        newSelectedRows.forEach((row, newIndex) => {
          newSelectedRowsMap[row.process_uuid] = newIndex;
        });
        return { selectedRows: newSelectedRows, selectedRowsMap: newSelectedRowsMap };
      }
      return state;
    }),
  addSelectedRow: (row) =>
    set((state) => {
      if (!state.selectedRowsMap.hasOwnProperty(row.process_uuid)) {
        const newSelectedRows = [...state.selectedRows, row];
        const newSelectedRowsMap = { ...state.selectedRowsMap, [row.process_uuid]: newSelectedRows.length - 1 };
        return { selectedRows: newSelectedRows, selectedRowsMap: newSelectedRowsMap };
      }
      return state;
    }),
}));
