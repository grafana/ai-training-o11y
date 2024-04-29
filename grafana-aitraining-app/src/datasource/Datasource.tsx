import { DataQueryRequest, DataQueryResponse, DataQueryResponseData, LoadingState, TestDataSourceResponse, arrayToDataFrame } from "@grafana/data";
import { getBackendSrv } from "@grafana/runtime";
import { RuntimeDataSource } from "@grafana/scenes";
import { DataQuery } from "@grafana/schema";
import { Observable, lastValueFrom } from "rxjs";

export const getProjects = (apiPrefix: string): Promise<DataQueryResponseData> => {
  return lastValueFrom(getBackendSrv().fetch({
    url: `${apiPrefix}/resources/projects`,
    method: 'GET',
  })).then((response) => {
    if (!response.ok) {
      throw response.data;
    }
    return response?.data;
  })
};


export class TrainingApiDatasource extends RuntimeDataSource {
  private apiPrefix: string;

  constructor(datasourceId: string, datasourceUid: string, pluginId: string) {
    super(datasourceId, datasourceUid);
    this.apiPrefix = `/api/plugins/${pluginId}`;
  }

  query(request: DataQueryRequest<DataQuery>): Promise<DataQueryResponse> | Observable<DataQueryResponse> {
    return getProjects(this.apiPrefix).then((projects) => {
      return {
        data: [
          arrayToDataFrame(projects),
        ],
        state: LoadingState.Done,
      };
    });
  }

  testDatasource(): Promise<TestDataSourceResponse> {
    return Promise.resolve({ status: 'success', message: 'Datasource frontend component works' });
  }
}
