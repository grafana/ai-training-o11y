import React from 'react';
import { RowData, QueryStatus } from 'utils/state';

interface TableTabProps {
  processQueryStatus: QueryStatus;
  rows: RowData[];
  isSelected: boolean[];
  setIsSelected: (index: number, value: boolean) => void;
  addSelectedRow: (row: RowData) => void;
  removeSelectedRow: (processUuid: string) => void;
}

export const TableTab: React.FC<TableTabProps> = ({
  processQueryStatus,
  rows,
  isSelected,
  setIsSelected,
  addSelectedRow,
  removeSelectedRow,
}: TableTabProps) => {
  // Get the unique column names from the data
  const columnNames = Array.from(new Set(rows.flatMap(Object.keys)));

  return (
    <div>
      <h2>Table</h2>
      {processQueryStatus === 'success' || <div>processQueryStatus</div>}
      <table style={styles.table}>
        <thead>
          <tr>
            <th style={styles.th}></th> {/* Empty header cell for multiselect */}
            {columnNames.map((column, index) => (
              <th key={index} style={styles.th}>
                {column}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((item, index) => (
            <tr key={index} style={styles.tr}>
              <td style={styles.td}>
                <input
                  type="checkbox"
                  checked={isSelected[index]}
                  onChange={(e) => {
                    if (e.target.checked) {
                      addSelectedRow(item);
                      setIsSelected(index, true);
                    } else {
                      removeSelectedRow(item.process_uuid);
                      setIsSelected(index, false);
                    }
                  }}
                />
              </td>
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
    borderCollapse: 'collapse' as 'collapse',
    width: '100%',
  },
  th: {
    backgroundColor: 'gray',
    border: '1px solid black',
    padding: '10px',
    textAlign: 'center' as 'center',
  },
  tr: {
    borderBottom: '1px solid #ddd',
  },
  td: {
    border: '1px solid #ddd',
    padding: '8px',
  },
  statusMessage: {
    padding: '10px',
    marginBottom: '10px',
    backgroundColor: '#e6f3ff',
    border: '1px solid #b8daff',
    borderRadius: '4px',
  },
  errorMessage: {
    padding: '10px',
    marginBottom: '10px',
    backgroundColor: '#f8d7da',
    border: '1px solid #f5c6cb',
    borderRadius: '4px',
    color: '#721c24',
  },
};

export default TableTab;
