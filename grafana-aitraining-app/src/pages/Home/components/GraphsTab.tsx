import React from 'react';
import { RowData } from 'utils/state';

interface GraphsProps {
  rows: RowData[];
}

export const GraphsTab: React.FC<GraphsProps> = ({ rows }) => {
  return (
    <div>
      <h2>Graphs Props</h2>
      <pre>{JSON.stringify(rows, null, 2)}</pre>
    </div>
  );
};
