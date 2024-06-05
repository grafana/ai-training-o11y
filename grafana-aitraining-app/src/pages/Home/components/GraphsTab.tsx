import React from 'react';

import { useProcessQueries } from 'hooks/useProcessQueries';
import { useTrainingAppStore, RowData } from 'utils/state';

interface GraphsProps {
  rows: RowData[];
}

export const GraphsTab: React.FC<GraphsProps> = ({ rows }) => {
  const { queryStatus, queryData, resetResults } = useTrainingAppStore();
  const { isReady, runQueries } = useProcessQueries();
  
  if (!isReady) {
    return <div>Loading...</div>;
  }

  if (queryStatus === 'running') {
    return <div>Running...

<button onClick={resetResults}>Reset Results</button>

    </div>;
  }

  if (queryStatus === 'empty') {
    runQueries();
  }

  console.log('query results', queryData);

  // At this point, queryData should be an object with keys being process_uuids
  // and all values should be present for use

  return (
    <div>

      <button onClick={resetResults}>Reset Results</button>

      <div style={{ marginBottom: '20px'}}>
        {Object.keys(queryData).map((key) => (
          <>
          Results for process: {key}
          <pre>{JSON.stringify(queryData[key].lokiData?.series[0].fields, null, 2)}</pre>
          </>
        ))}
      </div>


      <pre>
        JSON FOR TESTING BELOW:
        {JSON.stringify(rows, null, 2)}
      </pre>
    </div>
  );
};
