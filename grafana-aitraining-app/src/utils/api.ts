import { DataQueryResponseData } from "@grafana/data";
import { getBackendSrv } from "@grafana/runtime";
import { lastValueFrom } from "rxjs";

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

// plugin id set by: this.props.meta.id from the app.tsx

export function makeProcessGetter(pluginId: string) {
    return makeGetRequester(`/api/plugins/${pluginId}/resources/metadata/api/v1/processes`);
}
