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
    let lastStateTime = performance.now();

    const unsubscribe = () => {
      if (subRef !== null) {
        console.log('unsubscribing');
        subRef.unsubscribe();
      }
    }

    const runner = createQueryRunner();

    const startTime = performance.now();
    subRef = runner.get().subscribe({
      next: (data) => {
        const currentTime = performance.now();
        const timeSinceStart = currentTime - startTime;
        const timeSinceLastState = currentTime - lastStateTime;
        
        console.log('listening', data?.state, 
          `Time since start: ${timeSinceStart.toFixed(2)}ms`, 
          `Time since last state: ${timeSinceLastState.toFixed(2)}ms`
        );
        console.log(data);
        
        if (data?.state === 'Done') {
          console.log('done!!', 
            `Total time: ${timeSinceStart.toFixed(2)}ms`, 
            `Time since last state: ${timeSinceLastState.toFixed(2)}ms`
          );
          onResult(data);
          unsubscribe();
          resolve();
        }
        
        lastStateTime = currentTime;
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
