import { useCallback, useEffect, useState } from 'react';

import { PanelData, TimeRange } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

import { useCancelQuery, useQueryRunner } from './useQueryRunner';

/// Perform a set of data queries;
/// returning current panel along with isRefresh boolean, a refresh function, and a cancel function
export function useQueryResult(
  maxDataPoints: number,
  datasource: any
): [PanelData | undefined, (query: DataQuery, timeRange: TimeRange, timeZone: string) => void, () => void, boolean] {
  const runner = useQueryRunner();
  const cancelQuery = useCancelQuery();
  const [isRunning, setIsRunning] = useState<boolean>(false);

  const [value, update] = useState<PanelData | undefined>(undefined);
  const handleUpdate = useCallback(
    (v: PanelData | undefined) => {
      setIsRunning(false);
      update(v);
      // TODO: handle results
    },
    [update]
  );

  useEffect(() => {
    const s = runner.get().subscribe(handleUpdate);
    return () => s.unsubscribe();
  }, [runner, handleUpdate]);

  // todo: need a loading status
  // move time range/time zone here
  const onRunQuery = async (query: DataQuery, timeRange: TimeRange, timeZone: string) => {
    setIsRunning(true);
    await setTimeout(() => {}, 3000);
    const queries: DataQuery[] = [query];
    runner.run({
      timeRange,
      queries,
      datasource,
      maxDataPoints,
      minInterval: null,
      timezone: timeZone,
    });
  };

  return [value, onRunQuery, cancelQuery, isRunning];
}
