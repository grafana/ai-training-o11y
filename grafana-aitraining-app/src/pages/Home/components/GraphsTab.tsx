import React, { useEffect, useRef } from 'react';

import { useProcessQueries } from 'hooks/useProcessQueries';
import { useTrainingAppStore, RowData } from 'utils/state';
import { reshapeModelMetrics } from 'utils/reshapeModelMetrics';

interface GraphsProps {
  rows: RowData[];
}

export const GraphsTab: React.FC<GraphsProps> = ({ rows }) => {
  const {
    lokiQueryStatus,
    lokiQueryData,
    organizedLokiData,
    resetLokiResults,
    setOrganizedLokiData
  } = useTrainingAppStore();
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
        <button onClick={() => { resetLokiResults(); shouldRunQueries.current = true; }}>Reset Results</button>
      </div>
    );
  }

  if (!organizedLokiData) {
    return <div>No data</div>;
  }

  return (
    <div>
      <button onClick={() => { resetLokiResults(); shouldRunQueries.current = true; }}>Reset Results</button>
      
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
        <pre>
          {JSON.stringify(rows, null, 2)}
        </pre>
      </div>
    </div>
  );
};
