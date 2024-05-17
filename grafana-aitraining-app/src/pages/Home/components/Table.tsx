// Table.tsx
import React from 'react';

interface RowData {
  [key: string]: string | undefined;
}

interface TableProps {
  data: RowData[];
}

export const Table: React.FC<TableProps> = ({ data }) => {
  // Get the unique column names from the data
  const columnNames = Array.from(
    new Set(data.flatMap(Object.keys))
  );

  return (
    <div>
      <h2>Table</h2>
      <table>
        <thead>
          <tr>
            {columnNames.map((column, index) => (
              <th key={index}>{column}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((item, index) => (
            <tr key={index}>
              {columnNames.map((column, columnIndex) => (
                <td key={columnIndex}>{item[column]}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Table;
