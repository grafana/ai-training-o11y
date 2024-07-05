import React from 'react';

import { Pagination, Select, useStyles2 } from '@grafana/ui';

import { css } from '@emotion/css';

interface Option {
  label: string;
  value: number;
}

interface PagingProps {
  itemsPerPage: number;
  currentPage: number;
  totalItems: number;
  onItemsPerPageChange: (itemsPerPage: number) => void;
  onPageChange: (page: number) => void;
  options?: Option[];
  alignRight?: boolean;
}

const defaultOptions = [
  { label: '1', value: 1 },
  { label: '10', value: 10 },
  { label: '25', value: 25 },
  { label: '50', value: 50 },
  { label: '100', value: 100 },
];

export const Paging = ({
  itemsPerPage,
  currentPage,
  totalItems,
  onItemsPerPageChange,
  onPageChange,
  options = defaultOptions,
  alignRight = true,
}: PagingProps) => {
  const styles = useStyles2(getStyles);

  return (
    <div className={styles.container} style={alignRight ? { justifyContent: 'flex-end' } : {}}>
      <label>Items per page</label>
      <div className={styles.itemsPerPage}>
        <Select
          options={options}
          onChange={(val) => {
            if (val?.value !== undefined) {
              onItemsPerPageChange(val.value);
            }
          }}
          value={itemsPerPage}
        />
      </div>
      <Pagination
        currentPage={currentPage}
        numberOfPages={Math.ceil(totalItems / itemsPerPage)}
        onNavigate={(page) => onPageChange(page)}
        showSmallVersion
      />
    </div>
  );
};

const getStyles = () => ({
  container: css`
    display: flex;
    padding: 20px 0;
    & > label {
      width: 120px;
    }
  `,
  itemsPerPage: css`
    width: 70px;
  `,
});
