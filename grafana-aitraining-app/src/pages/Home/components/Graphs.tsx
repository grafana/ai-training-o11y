import React from 'react';

interface GraphsProps {
  // Define the props that the Graphs component expects
  data: any[];
}

export const Graphs: React.FC<GraphsProps> = ({ data }) => {
  return (
    <div>
      <h2>Graphs</h2>
      {/* Render the graphs using the data prop */}
      {/* Example: */}
      <ul>
        {data.map((item, index) => (
          <li key={index}>{item.value}</li>
        ))}
      </ul>
    </div>
  );
};
