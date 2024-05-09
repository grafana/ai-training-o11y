import { DataQueryRequest, DataQueryResponse, DataQueryResponseData, LoadingState, TestDataSourceResponse, arrayToDataFrame } from "@grafana/data";
import { getBackendSrv } from "@grafana/runtime";
import { RuntimeDataSource } from "@grafana/scenes";
import { DataQuery } from "@grafana/schema";
import { Observable, lastValueFrom } from "rxjs";

type JSONPrimitive = string | number | boolean | null;
type JSONObject = { [key: string]: JSONValue };
type JSONArray = JSONValue[];
type JSONValue = JSONPrimitive | JSONObject | JSONArray;

// Actually make a request from the plugin backend
function doRequest(fetchOptions: any): Promise<DataQueryResponseData> {
  return lastValueFrom(getBackendSrv().fetch(fetchOptions)).then((response) => {
    if (!response.ok) {
      throw response.data;
    }
    return response?.data;
  });
}

// Getter for endpoints without payload
export function makeGetRequester(apiPrefix: string): () => Promise<DataQueryResponseData> {
  return () => {
    const fetchOptions: any = {
      url: `${apiPrefix}`,
      method: 'GET',
    };

    return doRequest(fetchOptions);
  };
}

// Getter for endpoints with payload
export function makeGetRequesterWithPayload<T extends JSONObject>(apiPrefix: string): (payload: T) => Promise<DataQueryResponseData> {
  return (payload: T) => {
    const fetchOptions: any = {
      url: `${apiPrefix}`,
      method: 'GET',
      data: JSON.stringify(payload),
      headers: {
        'Content-Type': 'application/json',
      },
    };

    return doRequest(fetchOptions);
  };
}

export class TrainingApiDatasource extends RuntimeDataSource {
  private getProcesses: () => Promise<DataQueryResponseData>;

  constructor(datasourceId: string, datasourceUid: string, pluginId: string) {
    super(datasourceId, datasourceUid);
    this.getProcesses = makeGetRequester(`/api/plugins/${pluginId}/resources/metadata/api/v1/processes`);
  }

  query(request: DataQueryRequest<DataQuery>): Promise<DataQueryResponse> | Observable<DataQueryResponse> {
    return this.getProcesses().then((projects) => {
      return {
        data: [
          arrayToDataFrame(projects.data),
        ],
        state: LoadingState.Done,
      };
    });
  }

  testDatasource(): Promise<TestDataSourceResponse> {
    return Promise.resolve({ status: 'success', message: 'Datasource frontend component works' });
  }
}
