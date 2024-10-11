import React, { useEffect } from 'react';
import { RowData, useTrainingAppStore } from 'utils/state';
import { SceneGraph } from './SceneGraph';
import { PanelData, LoadingState, dateTime, TimeRange, DataFrame } from '@grafana/data';
import { ControlledCollapse } from '@grafana/ui';
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
  const { organizedData, appendToOrganizedData } = useTrainingAppStore();
  const [ loading, setLoading ] = React.useState<LoadingState>(LoadingState.NotStarted);

  useEffect(() => {
    if (rows.length > 0) {
      setLoading(LoadingState.Loading);
      const promises = rows.map((row) => {
        return getModelMetrics(row.process_uuid).then((response) => {
          response.data.forEach((item: any) => {
            const metricName = item.MetricName;
            const stepName = item.StepName;
            const section = metricName.includes('/') ? metricName.split('/')[0] : "general";

            const itemWithConfig: any = {
              ...item,
              fields: item.fields.map((field: any) => ({
                ...field,
                length: field.values.length,
                config: field.name === metricName ? {displayName: row.process_uuid } : {}
              }))
            };

            appendToOrganizedData(section, metricName, stepName, itemWithConfig);
          });
        });
      });

      Promise.all(promises).then(() => {
        setLoading(LoadingState.Done);
      });
    }
  }, [rows, getModelMetrics, appendToOrganizedData]);

  const startTime = dateTime();
  const endTime = dateTime();

  const tmpTimeRange: TimeRange = {
    from: startTime,
    to: endTime,
    raw: {
      from: startTime.toISOString(),
      to: endTime.toISOString(),
    },
  };

  const createPanelList = (section: string): MetricPanel[] => {
    console.log(`Creating panel list for section: ${section}`);
  
    if (!organizedData || !organizedData[section]) {
      console.log(`No data for section ${section}`);
      return [];
    }
  
    return Object.entries(organizedData[section]).map(([metricName, metricData]: any) => {
      const series: DataFrame[] = Object.entries(metricData).map(([stepName, stepData]: any): DataFrame => {
        const length: number = stepData.map((item: any) => item.fields[0].length).reduce((a: number, b: number) => Math.max(a, b), 0);
        return {
          name: metricName,
          fields: stepData[0].fields,
          length,
        }
      });

      const data: PanelData = {
        state: loading,
        timeRange: tmpTimeRange,
        series: series,
      }
  
      return {
        pluginId: 'xychart',
        title: metricName,
        data,
      }
    });
  };

  return (
    <div style={{ marginTop: '10px' }}>
      {organizedData && Object.keys(organizedData).map((section) => {
        const panels = createPanelList(section);
        return (
          <ControlledCollapse
            key={section}
            isOpen={true}
            label={`${section}`}
          >
            <SceneGraph panels={panels} />
          </ControlledCollapse>
        );
      })}
    </div>
  );
};
