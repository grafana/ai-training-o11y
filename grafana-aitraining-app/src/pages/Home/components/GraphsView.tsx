import React, { useEffect } from 'react';
import { RowData } from 'utils/state';
import { SceneGraph } from './SceneGraph';
import { PanelData, LoadingState, dateTime, TimeRange, FieldType } from '@grafana/data';
// import { ControlledCollapse } from '@grafana/ui';
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
  const getModelMetrics = useGetModelMetrics();

  useEffect(() => {
    if (rows.length > 0) {
      const organizedData: any = {};
      rows.forEach((row) => {
        const processUuid = row.process_uuid;
        getModelMetrics(processUuid).then((response) => {
          const data = response.data;
          console.log(data);
          // data should be an array, we want to loop over it
          for (const item of data) {
            const metricName = item?.MetricName;
            const stepName = item?.StepName;
            // If the MetricName has a slash in it, section is the first part of the MetricName
            // and key is the second part
            let section = "general";
            if (metricName.includes('/')) {
              section = metricName.split('/')[0];
            }

            // Init section if empty
            organizedData[section] = organizedData[section] ?? {};
            // Init metric if empty
            organizedData[section][metricName] = organizedData[section][metricName] ?? {};
            // Init step if empty
            organizedData[section][metricName][stepName] = organizedData[section][metricName][stepName] ?? [];
            // Append the data
            organizedData[section][metricName][stepName].push(item);

            console.log("Organized Data");
            console.log(organizedData);
          }
        });
      });
    }
  }, [rows, getModelMetrics]);

  interface Dataframe {
    // eslint-disable-next-line @typescript-eslint/array-type
    fields: {
      name: string;
      type: string;
      values: number[];
    }[];
  }
  
  interface Panel {
    metricName: string;
    dataframes: Dataframe[];
  }
  
  interface Section {
    sectionName: string;
    panels: Panel[];
  }
  
  const testData: Section[] = ['train', 'test'].map((sectionName: string) => {
    return {
      sectionName,
      panels: [
        {
          metricName: `${sectionName}/loss`,
          dataframes: [
            {
              fields: [
                {
                  name: 'step',
                  type: 'number',
                  values: [1, 2, 3],
                },
                {
                  name: 'value',
                  type: 'number',
                  values: [0.1, 0.2, 0.3],
                }
              ]
            }
          ],
        },
        {
          metricName: `${sectionName}/acc`,
          dataframes: [
            {
              fields: [
                {
                  name: 'step',
                  type: 'number',
                  values: [1, 2, 3],
                },
                {
                  name: 'value',
                  type: 'number',
                  values: [0.7, 0.8, 0.9], // Changed values for accuracy
                }
              ]
            }
          ],
        },
      ]
    };
  });

  console.log("testData");
  console.log(testData);

  const createPanelList = (section: Section): MetricPanel[] => {
    const dummyStart = dateTime(new Date('2021-09-04T00:00:00Z'));
    const dummyEnd = dateTime(new Date('2021-09-04T00:10:00Z'));
  
    const dummyTimeRange: TimeRange = {
      from: dummyStart,
      to: dummyEnd,
      raw: {
        from: dummyStart.toISOString(),
        to: dummyEnd.toISOString(),
      },
    };
    return section.panels.map((data: any) => {
      const panel: PanelData = {
        state: LoadingState.Done,
        timeRange: dummyTimeRange,
        series: data.dataframes,
      };
      return {
        pluginId: 'trend',
        title: data.metricName,
        data: panel,
      };
    });
  };

  const testSection: Section = {
    sectionName: "testSection",
    panels: [
      {
        metricName: "test/metric",
        dataframes: [
          {
            fields: [
              {
                name: "metric",
                type: "number",
                values: [1, 2, 3]
              },
              {
                name: "value",
                type: "number",
                values: [10, 20, 30]
              }
            ]
          }
        ]
      }
    ]
  };
  
  // Test the function
  console.log("testSection");
  const result = createPanelList(testSection);
  console.log(JSON.stringify(result, null, 2));

  const testPanel = {
    title: 'Test Panel',
    pluginId: 'xychart',
    data: {
      state: LoadingState.Done,
      series: [
        {
          fields: [
            { 
              name: 'value', 
              type: FieldType.number,
              values: [1, 2, 3],
              config: {}
            },
            { 
              name: 'value', 
              type: FieldType.number,
              values: [10, 20, 30],
              config: {}
            }
          ],
          length: 3 
        }
      ],
      timeRange: {
        from: dateTime(new Date('2021-09-04T00:00:00Z')),
        to: dateTime(new Date('2021-09-04T00:10:00Z')),
        raw: {
          from: '2021-09-04T00:00:00Z',
          to: '2021-09-04T00:10:00Z',
        }
      }
    }
  }

  console.log("testPanel");
  console.log(testPanel);

return (
  <div>
    <SceneGraph panels={[testPanel]} />
  </div>
);
}
