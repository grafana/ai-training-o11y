import React from 'react';
import { AppRootProps } from '@grafana/data';
import { sceneUtils } from '@grafana/scenes';
import { PluginPropsContext } from '../../utils/utils.plugin';
import { Routes } from '../Routes';
import { TrainingApiDatasource } from '../../datasource/Datasource';
import { QueryRunnerProvider } from 'hooks/useQueryRunner';

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

  render() {
    return (
      <QueryRunnerProvider>
        <PluginPropsContext.Provider value={this.props}>
          <Routes />
        </PluginPropsContext.Provider>
      </QueryRunnerProvider>
    );
  }
}
