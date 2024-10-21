import React, { useEffect } from 'react';
import { PanelData, LoadingState, DataFrame, dateTime, TimeRange } from '@grafana/data';
import { ControlledCollapse } from '@grafana/ui';
import { config } from '@grafana/runtime';

import { SceneGraph } from './SceneGraph';
import { RowData } from 'utils/state';

import { useGetModelMetrics } from 'utils/utils.plugin';

const palette = config.theme2.visualization.palette;

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
  const [ colors, setColors ] = React.useState<Map<string, string>>(new Map());

  useEffect(() => {
    if (rows.length > 0) {
      setLoading(LoadingState.Loading);
      const rowUUIDs = rows.map((row) => row.process_uuid);
      let newColors: Map<string, string> = new Map();
      rowUUIDs.forEach((uuid, i) => {
        newColors.set(uuid, palette[i % palette.length]);
      });
      setColors(newColors);
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
    const unit = panelData.series[0].name;
    const fields = panelData.series.map((s: any): DataFrame[] => {
      s.values = [
        ...s.values.map((v: string | undefined) => {
          if (v === undefined || v === null) {
            return undefined
          }
          return parseFloat(v)
      })
      ]
      return {
        ...s,
        config: {
          color: {
            mode: 'fixed',
            fixedColor: colors.get(s.name) || palette[0],
          },
        }
      }
    });
    return {
      pluginId: 'trend',
      title: `${panelData.title} per ${unit}`,
      data: {
        state: loading,
        timeRange: dummyTimeRange,
        series: [{
          fields,
          length: fields.length
        }],
      },
    }
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
