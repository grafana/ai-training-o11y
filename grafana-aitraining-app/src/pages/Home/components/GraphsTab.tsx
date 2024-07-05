import React, { useEffect, useRef } from 'react';

import { useProcessQueries } from 'hooks/useProcessQueries';
import { useTrainingAppStore, RowData } from 'utils/state';
import { reshapeModelMetrics } from 'utils/reshapeModelMetrics';
import { SceneGraph } from './SceneGraph';

import { PanelData, LoadingState, dateTime, TimeRange, FieldType } from '@grafana/data';

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

  if (Object.keys(organizedLokiData).length === 0) {
    return <div>No data</div>;
  }

  const startTime = dateTime(organizedLokiData.startTime);
  const endTime = dateTime(organizedLokiData.endTime);

  const tmpTimeRange: TimeRange = {
    from: startTime,
    to: endTime,
    raw: {
      from: startTime.toISOString(),
      to: endTime.toISOString(),
    },
  };

  const panelList: MetricPanel[] = organizedLokiData.meta.keys.map((key: string) => {
    const trendLine = Object.keys(organizedLokiData.data[key]).map((processId: any) => {
      const line = organizedLokiData.data[key][processId];
      return line.map((m: any, i: number) => i);
    })[0];

    const panel: PanelData = {
      state: LoadingState.Done,
      timeRange: tmpTimeRange,
      series: [
        {
          name: key,
          fields: [
            {
              name: 'Line',
              type: FieldType.number,
              values: trendLine,
              config: {},
            },
            ...Object.keys(organizedLokiData.data[key]).map((processId: any) => {
              const line = organizedLokiData.data[key][processId];
              return {
                name: 'Line',
                type: FieldType.number,
                values: [...line],
                config: {},
              };
            }),
          ],
          length: organizedLokiData.data[key].length,
        },
      ],
    };
    return {
      pluginId: 'trend',
      title: key,
      data: panel,
    };
  });

  return (
    <div style={{ marginTop: '10px' }}>
      <div>
        <SceneGraph panels={panelList} />
      </div>

      {/* remove this later, keeping here for debugging for now */}
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
