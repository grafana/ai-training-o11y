import React from 'react';

import { AppRootProps } from '@grafana/data';
import { sceneUtils } from '@grafana/scenes';

import { parse, stringify } from 'query-string';
import { QueryParamProvider } from 'use-query-params';
import { ReactRouter5Adapter } from 'use-query-params/adapters/react-router-5';

import { PluginPropsContext } from '../../utils/utils.plugin';
import { Routes } from '../Routes';
import { TrainingApiDatasource } from '../../datasource/Datasource';
import { doRequest, makeProcessGetter } from 'utils/api';

export class App extends React.PureComponent<AppRootProps> {
  componentDidMount() {
    try {
      sceneUtils.registerRuntimeDataSource({
        dataSource: new TrainingApiDatasource(
          'grafana-aitraining-app-datasource',
          'grafana-aitraining-app-datasource-uid',
          this.props.meta.id
        ),
      });
    } catch (e) {
      // eslint-disable-next-line no-console
      // Datasource already registered, probably
      console.error(e);
    }
  }

  getProcesses = makeProcessGetter(this.props.meta.id);

  getModelMetrics = (processUuid: string) => {
    const response = doRequest({
      url: `/api/plugins/${this.props.meta.id}/resources/metadata/api/v1/process/${processUuid}/model-metrics`,
      method: 'GET',
    });
    return response;
  }

  render() {
    return (
      <QueryParamProvider
        adapter={ReactRouter5Adapter}
        options={{
          searchStringToObject: parse,
          objectToSearchString: stringify,
        }}
      >
        <PluginPropsContext.Provider value={{ ...this.props, getProcesses: this.getProcesses, getModelMetrics: this.getModelMetrics }}>
          <Routes />
        </PluginPropsContext.Provider>
      </QueryParamProvider>
    );
  }
}
