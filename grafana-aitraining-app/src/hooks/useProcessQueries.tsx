import { useAsync } from 'react-use';
import { useTrainingAppStore } from 'utils/state';
import { runQuery } from 'utils/runQuery';
import { dateTime, TimeRange } from '@grafana/data';
import { getDataSourceSrv } from '@grafana/runtime';

const useProcessQueries = () => {
  const datasource = useAsync(async () => {
    return getDataSourceSrv().get('Loki');
  }, []);

  const isReady = datasource.loading !== true;

  const { selectedRows, appendResult, setQueryStatus } = useTrainingAppStore();

  const runQueries = async () => {
    console.log(`Starting all queries at ${new Date().toISOString()}`);
    setQueryStatus('loading');
  
    const queryPromises = selectedRows.map(async (processData, index) => {
      console.log(`Preparing query ${index} at ${new Date().toISOString()}`);
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
        expr: `{job="o11y"} | process_uuid = \`${processData.process_uuid}\` |= \`\``,
        queryType: 'range',
      };

      console.log(`Starting query ${index} execution at ${new Date().toISOString()}`);
      try {
        await runQuery({
          datasource: datasource.value,
          maxDataPoints: 100,
          queries: [query],
          timeRange: tmpTimeRange,
          timeZone: 'EST',
          onResult: (data: any) => {
            console.log(`Query ${index} completed at ${new Date().toISOString()}`);
            appendResult(processData, data);
          },
        });
        console.log(`Query ${index} promise resolved at ${new Date().toISOString()}`);
      } catch (error) {
        console.error(`Error in query ${index}:`, error);
        throw error; // Re-throw to be caught by Promise.all
      }
    });

    try {
      console.log(`Awaiting all queries at ${new Date().toISOString()}`);
      await Promise.all(queryPromises);
      console.log(`All queries completed at ${new Date().toISOString()}`);
      setQueryStatus('success');
    } catch (error) {
      console.error(`Error running queries at ${new Date().toISOString()}:`, error);
      setQueryStatus('error');
    }
  };

  return {
    isReady,
    runQueries,
  };
};

export { useProcessQueries };
