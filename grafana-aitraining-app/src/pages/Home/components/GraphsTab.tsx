import React, { useEffect, useRef } from 'react';
import { useProcessQueries } from 'hooks/useProcessQueries';
import { useTrainingAppStore, RowData } from 'utils/state';
import { reshapeModelMetrics } from 'utils/reshapeModelMetrics';
import { SceneGraph } from './SceneGraph';
import { PanelData, LoadingState, dateTime, TimeRange } from '@grafana/data';
import { ControlledCollapse } from '@grafana/ui';

export interface MetricPanel {
  pluginId: string;
  title: string;
  data: PanelData;
}

interface GraphsProps {
  rows: RowData[];
}

export const GraphsTab: React.FC<GraphsProps> = ({ rows }) => {
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
        <button
          onClick={() => {
            resetLokiResults();
            shouldRunQueries.current = true;
          }}
        >
          Reset Results
        </button>
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
