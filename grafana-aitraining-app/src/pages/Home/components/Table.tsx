// Table.tsx
import React from 'react';

interface TableProps {
  // Define the props that the Table component expects
  data: any[];
}

export const Table: React.FC<TableProps> = ({ data }) => {
  return (
    <div>
      <h2>Table</h2>
      {/* Render the table using the data prop */}
      {/* Example: */}
      <table>
        <thead>
          <tr>
            <th>Column 1</th>
            <th>Column 2</th>
          </tr>
        </thead>
        <tbody>
          {data.map((item, index) => (
            <tr key={index}>
              <td>{item.column1}</td>
              <td>{item.column2}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Table;
