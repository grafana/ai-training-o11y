import { TimeRange } from '@grafana/data';
import { createQueryRunner } from '@grafana/runtime';
import { DataQuery } from '@grafana/schema';

interface RunQueryProps {
  timeRange: TimeRange;
  queries: DataQuery[];
  datasource: any;
  maxDataPoints: number;
  minInterval?: string;
  timeZone: string;
  onResult: (data: any) => void;
}

export const runQuery = ({
  datasource,
  maxDataPoints,
  minInterval,
  queries,
  timeRange,
  timeZone,
  onResult
}: RunQueryProps): Promise<void> => {
  return new Promise((resolve, reject) => {
    let subRef: any = null;

    const unsubscribe = () => {
      if (subRef !== null) {
        subRef.unsubscribe();
      }
    }

    const runner = createQueryRunner();

    subRef = runner.get().subscribe({
      next: (data) => {        
        if (data?.state === 'Done') {
          onResult(data);
          unsubscribe();
          resolve();
        }
        
      },
      error: (error) => {
        console.error('Query error:', error);
        unsubscribe();
        reject(error);
      }
    });

    runner.run({
      timeRange,
      queries,
      datasource,
      maxDataPoints,
      minInterval,
      timezone: timeZone,
    });
  });
}
