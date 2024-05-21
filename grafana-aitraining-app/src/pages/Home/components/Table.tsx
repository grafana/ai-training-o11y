// Table.tsx
import React from 'react';
import { RowData } from 'utils/state';

interface TableProps {
  data: RowData[];
}

export const Table: React.FC<TableProps> = ({ data }) => {
  // Get the unique column names from the data
  const columnNames = Array.from(new Set(data.flatMap(Object.keys)));

  return (
    <div>
      <h2>Table</h2>
      <table style={styles.table}>
        <thead>
          <tr>
            {columnNames.map((column, index) => (
              <th key={index} style={styles.th}>
                {column}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((item, index) => (
            <tr key={index} style={styles.tr}>
              {columnNames.map((column, columnIndex) => (
                <td key={columnIndex} style={styles.td}>
                  {item[column]}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

const styles = {
  table: {
    borderCollapse: 'collapse' as 'collapse', // This is now a specific string value
    width: '100%',
  },
  th: {
    backgroundColor: 'gray',
    border: '1px solid black',
    padding: '10px',
    textAlign: 'center' as 'center', // This is now a specific string value
  },
  tr: {
    borderBottom: '1px solid #ddd',
  },
  td: {
    border: '1px solid #ddd',
    padding: '8px',
  },
};

export default Table;
