import { useTrainingAppStore } from 'utils/state';
import { runQuery } from 'utils/runQuery';
import { dateTime, TimeRange } from '@grafana/data';
import { FetchResponse, getBackendSrv, getDataSourceSrv } from '@grafana/runtime';
import { lastValueFrom } from 'rxjs';
import { useAsync } from 'react-use';

export const getSettings = async (pluginId: string) => {
  const response = getBackendSrv().fetch({
    url: `/api/plugins/${pluginId}/settings`,
    method: 'get',
  });

  const dataResponse = await lastValueFrom(response);
  const { lokiDatasourceName, mimirDatasourceName, metadataUrl } = (dataResponse as FetchResponse<any>).data.jsonData;
  return { lokiDatasourceName, mimirDatasourceName, metadataUrl };
}

const useProcessQueries = () => {
  let datasource: any = null;

  // Define the async function to get the datasource
  const fetchDatasource = async () => {
    try {
      const settings = await getSettings('grafana-aitraining-app');
      return getDataSourceSrv().get(settings.lokiDatasourceName);
    } catch (error) {
      console.error('Error getting datasource settings:', error);
      throw error;
    }
  };
  
  // Use useAsync at the top level
  datasource = useAsync(fetchDatasource, []);

  const isReady = datasource.loading !== true;

  const { selectedRows, appendLokiResult: appendResult, setLokiQueryStatus } = useTrainingAppStore();

  const runQueries = async () => {
    setLokiQueryStatus('loading');
  
    const queryPromises = selectedRows.map(async (processData, index) => {
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

      try {
        await runQuery({
          datasource: datasource.value,
          maxDataPoints: 100,
          queries: [query],
          timeRange: tmpTimeRange,
          timeZone: 'EST',
          onResult: (data: any) => {
            appendResult(processData, data);
          },
        });
      } catch (error) {
        console.error(`Error in query ${index}:`, error);
        throw error; // Re-throw to be caught by Promise.all
      }
    });

    try {
      await Promise.all(queryPromises);
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
