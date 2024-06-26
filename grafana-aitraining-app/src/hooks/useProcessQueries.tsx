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

  const { selectedRows, appendLokiResult: appendResult, setLokiQueryStatus } = useTrainingAppStore();

  const runQueries = async () => {
    console.log(`Starting all queries at ${new Date().toISOString()}`);
    setLokiQueryStatus('loading');
  
    const queryPromises = selectedRows.map(async (processData, index) => {
      console.log(`Preparing query ${index} at ${new Date().toISOString()}`);

      const startDate = dateTime(processData.start_time);
      // If the process is still running, use the current time as the end time
      let endDate = processData.status === 'running' ? dateTime(new Date()): dateTime(processData.end_time);

      const tmpTimeRange: TimeRange = {
        from: startDate,
        to: endDate,
        raw: {
          from: startDate.toISOString(),
          to: endDate.toISOString()
        },
      };
  
      const query = {
        refId: 'A',
        expr: `{job="o11y"} | process_uuid = \`${processData.process_uuid}\``,
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
      setLokiQueryStatus('success');
    } catch (error) {
      console.error(`Error running queries at ${new Date().toISOString()}:`, error);
      setLokiQueryStatus('error');
    }
  };

  return {
    isReady,
    runQueries,
  };
};

export { useProcessQueries };
