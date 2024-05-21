import { useCallback, useEffect, useState } from 'react';

import { LoadingState, PanelData, TimeRange } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

import { useCancelQuery, useQueryRunner } from './useQueryRunner';

/// Perform a set of data queries;
/// returning current panel along with isRefresh boolean, a refresh function, and a cancel function
export function useQueryResult(
  maxDataPoints: number,
  timeRange: TimeRange,
  timeZone: string,
  datasource: any,
): [PanelData | undefined, (query: DataQuery) => void, () => void] {
  const runner = useQueryRunner();
  const cancelQuery = useCancelQuery();

  const [value, update] = useState<PanelData | undefined>(undefined);
  const handleUpdate = useCallback(
    (v: PanelData | undefined) => {
      update(v);
    },
    [update]
  );

  useEffect(() => {
    const s = runner.get().subscribe(handleUpdate);
    return () => s.unsubscribe();
  }, [runner, handleUpdate]);

  const onRunQuery = async (query: DataQuery) => {
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

  return [value, onRunQuery, cancelQuery];
}
