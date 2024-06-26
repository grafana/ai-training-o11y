import React, { useContext } from 'react';
import { AppRootProps, DataQueryResponseData } from '@grafana/data';

export interface PluginProps extends AppRootProps {
  getProcesses: () => Promise<DataQueryResponseData>;
}

export const PluginPropsContext = React.createContext<PluginProps | null>(null);

export const usePluginProps = () => {
  const pluginProps = useContext(PluginPropsContext);
  if (!pluginProps) {
    throw new Error('usePluginProps must be used within a PluginPropsProvider');
  }
  return pluginProps;
};

export const useGetProcesses = () => {
  const pluginProps = usePluginProps();
  return pluginProps.getProcesses;
};
