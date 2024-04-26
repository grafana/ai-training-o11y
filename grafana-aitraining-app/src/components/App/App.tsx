import React from 'react';
import { AppRootProps } from '@grafana/data';
import { sceneUtils } from '@grafana/scenes';
import { PluginPropsContext } from '../../utils/utils.plugin';
import { Routes } from '../Routes';
import { TrainingApiDatasource } from '../../datasource/Datasource';



export class App extends React.PureComponent<AppRootProps> {
  componentDidMount() {
    sceneUtils.registerRuntimeDataSource({ dataSource: new TrainingApiDatasource('grafana-aitraining-app-datasource', 'grafana-aitraining-app-datasource-uid', this.props.meta.id)});
  }
  
  render() {
    return (
      <PluginPropsContext.Provider value={this.props}>
        <Routes />
      </PluginPropsContext.Provider>
    );
  }
}
