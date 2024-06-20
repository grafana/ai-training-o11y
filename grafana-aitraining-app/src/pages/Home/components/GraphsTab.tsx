import React, { useEffect } from 'react';

import { useProcessQueries } from 'hooks/useProcessQueries';
import { useTrainingAppStore, RowData } from 'utils/state';
import { reshapeModelMetrics } from 'utils/reshapeModelMetrics';

interface GraphsProps {
  rows: RowData[];
}

export const GraphsTab: React.FC<GraphsProps> = ({ rows }) => {
  const { queryStatus, queryData, organizedData, resetResults, setOrganizedData } = useTrainingAppStore();
  const { isReady, runQueries } = useProcessQueries();

  useEffect(() => {
    console.log('useEffect firing', {queryStatus, queryData});
    if (queryStatus === 'success' && Object.keys(queryData).length > 0) {
      const organized = reshapeModelMetrics(queryData);
      setOrganizedData(organized);
      console.log('Organized Data:', organized);
    }
  }, [queryStatus, queryData, setOrganizedData]);

  useEffect(() => {
    if (organizedData) {
      console.log('Organized Data from own effect:', organizedData);
    }
  }, [organizedData]);

  // Log queryData if it is not undefined
  useEffect(() => {
    if (queryData) {
      console.log('Query Data:', queryData);
    }
  }, [queryData]);

  if (!isReady) {
    return <div>Loading...</div>;
  }

  if (queryStatus === 'loading') {
    return (
      <div>
        Running...
        <button onClick={resetResults}>Reset Results</button>
      </div>
    );
  }

  if (queryStatus === 'idle') {
    runQueries();
  }

  if (!queryData) {
    return <div>No data</div>;
  }

  return (
    <div>
      <button onClick={resetResults}>Reset Results</button>
      <div style={{ marginBottom: '20px' }}>

      </div>

      <div style={{ marginBottom: '20px' }}>
        {Object.keys(queryData).map((key) => (
          <React.Fragment key={key}>
            Results for process: {key}
            <pre>{JSON.stringify(queryData[key].lokiData?.series[0].fields, null, 2)}</pre>
          </React.Fragment>
        ))}
      </div>

      <pre>
        JSON FOR TESTING BELOW:
        {JSON.stringify(rows, null, 2)}
      </pre>
    </div>
  );
};
