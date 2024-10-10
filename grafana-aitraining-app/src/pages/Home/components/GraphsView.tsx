import React, { useEffect } from 'react';
import { RowData } from 'utils/state';
import { SceneGraph } from './SceneGraph';
import { PanelData, LoadingState, dateTime, TimeRange } from '@grafana/data';
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
  
    return Object.entries(organizedData[section]).map(([key, metricData]: any) => {
      console.log(`Processing metric: ${key}`);
  
      const stepDataKey = Object.keys(metricData)[0];
      const stepData = metricData[stepDataKey];
  
      if (!stepData || !stepData[0] || !stepData[0].fields) {
        console.log(`Invalid step data for metric ${key}`);
        return {
          pluginId: 'xyplot',
          title: key,
          data: {
            state: LoadingState.Error,
            series: [],
            timeRange: tmpTimeRange,
          },
        };
      }
  
      const fieldsLength = stepData[0].fields[0].values.length;
  
      const panel: PanelData = {
        state: LoadingState.Done,
        timeRange: tmpTimeRange,
        series: [{
          name: key,
          fields: stepData[0].fields,
          length: fieldsLength,
        }],
      };
  
      console.log(`Created panel for metric ${key}:`, JSON.stringify(panel, null, 2));
  
      return {
        pluginId: 'xychart',
        title: key,
        data: panel,
      };
    });
  };

  return (
    <div style={{ marginTop: '10px' }}>
      {organizedData && Object.keys(organizedData).map((section) => {
        console.log(`Rendering section: ${section}`);
        const panels = createPanelList(section);
        console.log(`Panels for section ${section}:`, panels);
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
