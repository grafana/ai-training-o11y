import React, { useEffect } from 'react';
import { RowData } from 'utils/state';
import { PanelData, LoadingState } from '@grafana/data';
import { useGetModelMetrics } from 'utils/utils.plugin';

export interface MetricPanel {
  pluginId: string;
  title: string;
  data: PanelData;
}

interface GraphsProps {
  rows: RowData[];
}

export const GraphsView: React.FC<GraphsProps> = ({ rows }) => {
  // WIP:
  const getModelMetrics = useGetModelMetrics();
  const [ metrics, setMetrics ] = React.useState<any>();
  const [ loading, setLoading ] = React.useState<LoadingState>(LoadingState.NotStarted);

  useEffect(() => {
    if (rows.length > 0) {
      setLoading(LoadingState.Loading);
      const rowUUIDs = rows.map((row) => row.process_uuid);
      console.log("uuids");
      console.log(rowUUIDs);
      getModelMetrics(rowUUIDs).then((metrics) => {
        console.log("metrics");
        console.log(metrics);
        setMetrics(metrics);
        setLoading(LoadingState.Done);
      });
    }
  }, [rows, getModelMetrics]);

  console.log('ignore');
  console.log(loading, metrics);

  return (
    <div >
      This loaded
    </div>
  );
};
