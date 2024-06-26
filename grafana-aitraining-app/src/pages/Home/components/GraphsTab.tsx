import React, { useEffect, useRef } from 'react';

import { useProcessQueries } from 'hooks/useProcessQueries';
import { useTrainingAppStore, RowData } from 'utils/state';
import { reshapeModelMetrics } from 'utils/reshapeModelMetrics';
import { SceneGraph } from './SceneGraph';

import { PanelData, DataFrame, LoadingState, dateTime, TimeRange, FieldType } from '@grafana/data';

interface GraphsProps {
  rows: RowData[];
}

export const GraphsTab: React.FC<GraphsProps> = ({ rows }) => {
  const { queryStatus, queryData, organizedData, resetResults, setOrganizedData } = useTrainingAppStore();
  const { isReady, runQueries } = useProcessQueries();
  const shouldRunQueries = useRef(true);

  useEffect(() => {
    if (isReady && rows.length > 0 && shouldRunQueries.current) {
      shouldRunQueries.current = false;
      resetResults();
      runQueries();
    }
  }, [isReady, rows, resetResults, runQueries]);

  useEffect(() => {
    shouldRunQueries.current = true;
  }, [rows]);

  useEffect(() => {
    if (queryStatus === 'success' && Object.keys(queryData).length > 0) {
      const organized = reshapeModelMetrics(queryData);
      setOrganizedData(organized);
    }
  }, [queryStatus, queryData, setOrganizedData]);

  if (!isReady) {
    return <div>Loading...</div>;
  }

  if (queryStatus === 'loading') {
    return (
      <div>
        Running...
        <button
          onClick={() => {
            resetResults();
            shouldRunQueries.current = true;
          }}
        >
          Reset Results
        </button>
      </div>
    );
  }

  if (!organizedData) {
    return <div>No data</div>;
  }

  console.log('final data', { queryData, organizedData });

  /// ---- FAKE DATA BELOW -----

  /// FAKE TIME RANGE
  const tmpTimeRange: TimeRange = {
    from: dateTime('2024-06-26T00:01:00.001Z'), // startDate,
    to: dateTime('2024-06-26T10:30:00.001Z'), // endDate,
    raw: {
      from: dateTime('2024-06-26T00:01:00.001Z').toISOString(),
      to: dateTime('2024-06-26T10:30:00.001Z').toISOString()
    },
  };

  /// FAKE DATA FRAME
    const sampleFrame: DataFrame = {
      name: '',
      fields: [
        { name: 'Time', type: FieldType.time, values: [1,2,3,4,5,6,7,8], config: {} },
        { name: 'Time', type: FieldType.number, values: [1,2,3,4,5,6,7,8], config: {} },
      ],
      length: 2,
    };

  /// FAKE DATA PANE
  const samplePanel: PanelData = {
    state: LoadingState.Done,
    series: [
      sampleFrame
    ],
    timeRange: tmpTimeRange,
  };


  return (
    <div>
      <div style={{ marginBottom: '20px' }}>
        <h3>Scenes:</h3>
        <SceneGraph data={samplePanel} />
      </div>

      <button
        onClick={() => {
          resetResults();
          shouldRunQueries.current = true;
        }}
      >
        Reset Results
      </button>

      <div style={{ marginBottom: '20px' }}>
        <h3>Organized Data:</h3>
        {organizedData ? <pre>{JSON.stringify(organizedData, null, 2)}</pre> : <p>No organized data available</p>}
      </div>

      <div style={{ marginBottom: '20px' }}>
        <h3>Query Data:</h3>
        {Object.keys(queryData).map((key) => (
          <React.Fragment key={key}>
            <h4>Results for process: {key}</h4>
            <pre>{JSON.stringify(queryData[key].lokiData?.series[0].fields, null, 2)}</pre>
          </React.Fragment>
        ))}
      </div>

      <div>
        <h3>Selected Rows:</h3>
        <pre>{JSON.stringify(rows, null, 2)}</pre>
      </div>
    </div>
  );
};
