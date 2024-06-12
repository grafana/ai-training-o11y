import { useAsync } from 'react-use';

import { useTrainingAppStore } from 'utils/state';
import { runQuery } from 'utils/runQuery';

import { dateTime, TimeRange } from '@grafana/data';
import { getDataSourceSrv } from '@grafana/runtime';

// takes selected rows from state, and attempts to run a query for each
const useProcessQueries = () => {
  // get the datasource we need
  const datasource = useAsync(async () => {
    return getDataSourceSrv().get('Loki');
  }, []);

  const isReady = datasource.loading !== true;

  const { selectedRows, appendResult, setQueryStatus } = useTrainingAppStore();

  // export a function to run a query for each process in state
  const runQueries = () => {
    const doneCount = selectedRows.length;
    let currentCount = 0;

    setQueryStatus('loading');

    selectedRows.map((processData) => {
      // build the query here, including setting a time range and other details
      // using the processsData json values
      const startDate = new Date();
      startDate.setHours(startDate.getHours() - 7200);
      const endDate = dateTime(new Date());
      const tmpTimeRange: TimeRange = {
        from: dateTime(startDate),
        to: endDate,
        raw: { from: startDate.toLocaleString(), to: endDate.toLocaleString() },
      };

      const query = {
        refId: 'A',
        expr: '{job="o11y"}',
        queryType: 'range',
      };

      // run the query
      runQuery({
        datasource: datasource.value,
        maxDataPoints: 100,
        queries: [query],
        timeRange: tmpTimeRange,
        timeZone: 'EST',
        onResult: (data: any) => {
          currentCount++;
          // if this is the last process completed, mark the results as finished
          if (currentCount === doneCount) {
            setQueryStatus('success');
          }
          appendResult(processData, data);
        },
      });
    });
  };

  return {
    isReady,
    runQueries,
  };
};

export { useProcessQueries };
