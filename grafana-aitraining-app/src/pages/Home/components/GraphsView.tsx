import React, { useEffect, useRef } from 'react';
import { useProcessQueries } from 'hooks/useProcessQueries';
import { useTrainingAppStore, RowData } from 'utils/state';
import { reshapeModelMetrics } from 'utils/reshapeModelMetrics';
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
            organizedData[section][metricName][stepName].push({...item, config: {}});

            console.log("Organized Data");
            console.log(organizedData);
          }
        });
      });
    }
  }, [rows, getModelMetrics]);
  // End WIP


  const { lokiQueryStatus, lokiQueryData, organizedLokiData, resetLokiResults, setOrganizedLokiData } =
    useTrainingAppStore();
  const { isReady, runQueries } = useProcessQueries();
  const shouldRunQueries = useRef(true);

  useEffect(() => {
    if (isReady && rows.length > 0 && shouldRunQueries.current) {
      shouldRunQueries.current = false;
      resetLokiResults();
      runQueries();
    }
  }, [isReady, rows, resetLokiResults, runQueries]);

  useEffect(() => {
    shouldRunQueries.current = true;
  }, [rows]);

  useEffect(() => {
    if (lokiQueryStatus === 'success' && Object.keys(lokiQueryData).length > 0) {
      const organized = reshapeModelMetrics(lokiQueryData);
      setOrganizedLokiData(organized);
    }
  }, [lokiQueryStatus, lokiQueryData, setOrganizedLokiData]);

  if (!isReady) {
    return <div>Loading...</div>;
  }

  if (lokiQueryStatus === 'loading') {
    return (
      <div>
        Running...
      </div>
    );
  }

  if (!organizedLokiData || !organizedLokiData.data || Object.keys(organizedLokiData.data).length === 0) {
    return <div>No data available</div>;
  }

  const startTime = organizedLokiData.meta.startTime ? dateTime(organizedLokiData.meta.startTime) : dateTime();
  const endTime = organizedLokiData.meta.endTime ? dateTime(organizedLokiData.meta.endTime) : dateTime();

  const tmpTimeRange: TimeRange = {
    from: startTime,
    to: endTime,
    raw: {
      from: startTime.toISOString(),
      to: endTime.toISOString(),
    },
  };

  const createPanelList = (section: string): MetricPanel[] => {
    if (!organizedLokiData.meta.sections[section]) {
      return [];
    }
    return organizedLokiData.meta.sections[section].map((key: string) => {
      if (!organizedLokiData.data[section] || !organizedLokiData.data[section][key]) {
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
      const panel: PanelData = {
        state: LoadingState.Done,
        timeRange: tmpTimeRange,
        series: [organizedLokiData.data[section][key]],
      };
      return {
        pluginId: 'trend',
        title: key,
        data: panel,
      };
    });
  };

  return (
    <div style={{ marginTop: '10px' }}>
      {organizedLokiData.meta.sections && Object.entries(organizedLokiData.meta.sections).map(([section, keys]) => (
        <ControlledCollapse
          key={section}
          isOpen={true}
          label={`${section}`}
        >
          <SceneGraph panels={createPanelList(section)} />
        </ControlledCollapse>
      ))}

      {/* Debug section (hidden) */}
      <div style={{ display: 'none' }}>
        <button
          onClick={() => {
            resetLokiResults();
            shouldRunQueries.current = true;
          }}
        >
          Reset Results
        </button>

        <div style={{ marginBottom: '20px' }}>
          <h3>Organized Data:</h3>
          {organizedLokiData ? (
            <pre>{JSON.stringify(organizedLokiData, null, 2)}</pre>
          ) : (
            <p>No organized data available</p>
          )}
        </div>

        <div style={{ marginBottom: '20px' }}>
          <h3>Query Data:</h3>
          {Object.keys(lokiQueryData).map((key) => (
            <React.Fragment key={key}>
              <h4>Results for process: {key}</h4>
              <pre>{JSON.stringify(lokiQueryData[key].lokiData?.series[0].fields, null, 2)}</pre>
            </React.Fragment>
          ))}
        </div>

        <div>
          <h3>Selected Rows:</h3>
          <pre>{JSON.stringify(rows, null, 2)}</pre>
        </div>
      </div>
    </div>
  );
};
