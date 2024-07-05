import React from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { useStyles2, Checkbox } from '@grafana/ui';

import { css } from '@emotion/css';

import { RowData } from 'utils/state';

interface ProcessListItemProps {
  process: RowData;
  isSelected: boolean;
  addSelectedRow: (row: RowData) => void;
  removeSelectedRow: (processUuid: string) => void;
}

export const ProcessListItem: React.FC<ProcessListItemProps> = ({
  process,
  isSelected,
  addSelectedRow,
  removeSelectedRow,
}) => {
  const styles = useStyles2(getStyles);
  const startDate = new Date(process.start_time).toLocaleDateString();
  const startTime = new Date(process.start_time).toLocaleTimeString();
  const endTime = process.end_time !== undefined ? new Date(process.end_time).toLocaleTimeString() : '';

  return (
    <div key={process.process_uuid}>
      <div
        className={styles.card}
        onClick={() => {
          if (isSelected) {
            removeSelectedRow(process.process_uuid);
          } else {
            addSelectedRow(process);
          }
        }}
      >
        <div className={styles.checkboxWrapper}>
          <Checkbox
            value={isSelected}
          />
        </div>
        <div>{process.process_uuid}</div>
        <div className={styles.details}>
        {process.project} | {startDate}: {startTime} to {endTime} 
        </div>
      </div>
    </div>
  );
};


export const ProcessListItemPlaceholder: React.FC = () => {
  const styles = useStyles2(getStyles);

  return (
    <div key="placeholder">
      <div
        className={styles.card}
      >
        <div>
          There are no processes available.
        </div>
      </div>
    </div>
  );
};

const getStyles = (theme: GrafanaTheme2) => {
  return {
    empty: css``,

    card: css`
      display: flex;
      justify-content: flex-start;
      padding: 18px;
      background-color: ${theme.colors.background.secondary};
      margin-top: 8px;
      margin-bottom: 8px;
      border: 1px solid ${theme.colors.border.weak};
      border-radius: 2px;
      cursor: pointer;
    `,
    details: css`
      margin-left: auto;
    `,
    checkboxWrapper: css`
      margin-right: 8px;
    `,
  };
};

export default ProcessListItem;
