import React, { useEffect } from 'react';
import { RowData } from 'utils/state';
import { PanelData, LoadingState, DataFrame, dateTime, TimeRange, FieldType } from '@grafana/data';
import { useGetModelMetrics } from 'utils/utils.plugin';
import { ControlledCollapse } from '@grafana/ui';
import { SceneGraph } from './SceneGraph';
// import { SceneGraph } from './SceneGraph';

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
      getModelMetrics(rowUUIDs).then((metrics) => {
        if (metrics?.status !== "success") {
          setLoading(LoadingState.Error);
          return;
        } 
        setMetrics(metrics?.data?.sections);
        setLoading(LoadingState.Done);
      });
    }
  }, [rows, getModelMetrics]);

  function makePanelFromData(panelData: any) {
    const startTime = dateTime();
    const endTime = dateTime();
    const dummyTimeRange: TimeRange = {
      from: startTime,
      to: endTime,
      raw: {
        from: startTime.toISOString(),
        to: endTime.toISOString(),
      },
    };
    const fields = panelData.series.map((s: any): DataFrame[] => {
      return {
        ...s,
        values: s.type === FieldType.number && typeof(s.values[0]) === 'string' ? s.values.map((v: any) => parseFloat(v)) : s.values,
        config: {}
      }
    });
    const ret = {
      pluginId: 'trend',
      title: panelData.title,
      data: {
        state: loading,
        timeRange: dummyTimeRange,
        series: [{
          fields,
          length: fields.length
        }],
      },
    }
    return ret
  }

  return (
    <div style={{ marginTop: '10px' }}>
    {loading === LoadingState.Loading && <div>Loading...</div>}
    {metrics && Object.keys(metrics).map((section) => {
        return (
        <ControlledCollapse
            key={section}
            isOpen={true}
            label={`${section}`}
        >
            {<SceneGraph panels={metrics[section].map(makePanelFromData)} />}
        </ControlledCollapse>
        );
    })}
    </div>
  );
};
