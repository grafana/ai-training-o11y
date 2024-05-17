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
