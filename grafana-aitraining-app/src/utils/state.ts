import { create } from 'zustand'

interface TabState {
  tab: 'table' | 'graphs'
  set: (tab: string) => void
}

export const useTabStore = create<TabState>()((set) => ({
  tab: 'table',
  set: (to) => set((state) => ({ tab: to as 'table' | 'graphs'})),
}))
