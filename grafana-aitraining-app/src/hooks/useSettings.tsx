import { FetchResponse, getBackendSrv } from '@grafana/runtime';
import { lastValueFrom } from 'rxjs';
import { useAsync } from 'react-use';

export const getSettings = async (pluginId: string) => {
  const response = getBackendSrv().fetch({
    url: `/api/plugins/${pluginId}/settings`,
    method: 'get',
  });

  const dataResponse = await lastValueFrom(response);
  const { lokiDatasourceName, mimirDatasourceName, metadataUrl, stackId } = (dataResponse as FetchResponse<any>).data.jsonData;
  return { lokiDatasourceName, mimirDatasourceName, metadataUrl, stackId };
}

export const useSettings = () => {
  const fetchSettings = async () => {
    try {
      return await getSettings('grafana-aitraining-app');
    } catch (error) {
      console.error('Error getting datasource settings:', error);
      throw error;
    }
  };
  const settings = useAsync(fetchSettings, []);
  return {
    isReady: settings.loading !== true,
    settings: settings.value,
    error: settings.error,
  }
}
