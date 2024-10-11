import React, { useEffect } from 'react';
import { RowData } from 'utils/state';
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
  const [organizedData, setOrganizedData] = React.useState();

  useEffect(() => {
    if (rows.length > 0) {
      const newData: any = {};
      const promises = rows.map((row) => {
        return getModelMetrics(row.process_uuid).then((response) => {
          const data = response.data;
          for (const item of data) {
            const metricName = item?.MetricName;
            const stepName = item?.StepName;
            let section = metricName.includes('/') ? metricName.split('/')[0] : "general";
  
            newData[section] = newData[section] ?? {};
            newData[section][metricName] = newData[section][metricName] ?? {};
            newData[section][metricName][stepName] = newData[section][metricName][stepName] ?? [];
  
            const itemWithConfig = {
              ...item,
              fields: item.fields.map((field: any) => ({
                ...field,
                name: field.name === metricName ? row.process_uuid : field.name,
                length: field.values.length,
                config: {}
              }))
            };
            newData[section][metricName][stepName].push(itemWithConfig);
          }
        });
      });
  
      Promise.all(promises).then(() => {
        setOrganizedData(newData);
      });
    }
  }, [rows, getModelMetrics]);

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
        state: LoadingState.Done,
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
