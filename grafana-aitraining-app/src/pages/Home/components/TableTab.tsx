import React, { useEffect } from 'react';
import { RowData } from 'utils/state';

// import { Loki } from 'components/Datasources/Loki';

interface TableTabProps {
  rows: RowData[];
  isSelected: boolean[];
  setIsSelected: (index: number, value: boolean) => void;
  addSelectedRow: (row: RowData) => void;
  removeSelectedRow: (processUuid: string) => void;
}

export const TableTab: React.FC<TableTabProps> = ({
  rows,
  isSelected,
  setIsSelected,
  addSelectedRow,
  removeSelectedRow,
}: TableTabProps) => {


  useEffect(() => {
    console.log(isSelected);
  }, [isSelected]);

  // Get the unique column names from the data
  const columnNames = Array.from(new Set(rows.flatMap(Object.keys)));

  return (
    <div>
      <h2>Table</h2>
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

export default TableTab;
