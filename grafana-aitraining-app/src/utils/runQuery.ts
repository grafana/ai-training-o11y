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

export const runQuery = async ({
  datasource,
  maxDataPoints,
  minInterval,
  queries,
  timeRange,
  timeZone,
  onResult
}: RunQueryProps) => {
  let subRef: any = null;

  const unsubscribe = () => {
    if (subRef !== null) {
      console.log('unsubscribing');
      subRef.unsubscribe();
    }
  }

  const run = async () => {
    const runner = createQueryRunner();

    subRef = runner.get().subscribe((data) => {
      console.log('listening', data?.state);
      if (data?.state === 'Done') {
        console.log('done!!');
        onResult(data);
        unsubscribe();
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
  }

  run();
}
