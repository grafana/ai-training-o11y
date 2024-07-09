import React from 'react';

import { NumberParam, useQueryParam } from 'use-query-params';

import { ProcessListItem, ProcessListItemPlaceholder } from './ProcessListItem';
import { Paging } from 'components/Paging';
import { RowData, QueryStatus } from 'utils/state';

interface ProcessListProps {
  processQueryStatus: QueryStatus;
  rows: RowData[];
  selectedRows: RowData[];
  addSelectedRow: (row: RowData) => void;
  removeSelectedRow: (processUuid: string) => void;
}

export const ProcessList: React.FC<ProcessListProps> = ({
  processQueryStatus,
  rows,
  selectedRows,
  addSelectedRow,
  removeSelectedRow,
}) => {
  const [page, setPage] = useQueryParam('page', NumberParam);
  const [pageSize, setPageSize] = useQueryParam('pageSize', NumberParam);

  // Unauthorized, notFound and serverError unused due to odd behavior in api.ts.
  // Left for future resolution
  if (processQueryStatus === 'loading') {
    return <div>Loading...</div>;
  } else if (processQueryStatus === 'unauthorized') {
    return <div>Unauthorized. Please log in.</div>;
  } else if (processQueryStatus === 'notFound') {
    return <div>Resource not found.</div>;
  } else if (processQueryStatus === 'serverError') {
    return <div>Server error. Please try again later.</div>;
  } else if (processQueryStatus === 'error') {
    return <div>A server error occurred. Please try again.</div>;
  }

  const totalItems = rows?.length ?? 0;
  const currentPage = page ?? 1;
  const itemsPerPage = pageSize ?? 10;
  const pagedRows = rows?.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage) ?? [];

  return (
    <div>
      {pagedRows.map((item) => (
        <ProcessListItem
          key={item.process_uuid}
          process={item}
          isSelected={selectedRows.filter(r => r.process_uuid === item.process_uuid).length > 0}
          addSelectedRow={addSelectedRow}
          removeSelectedRow={removeSelectedRow}
        />
      ))}
      {pagedRows.length === 0 ? <ProcessListItemPlaceholder /> : null}
      {totalItems > 0 ? (
        <Paging
          currentPage={currentPage}
          itemsPerPage={itemsPerPage}
          totalItems={totalItems}
          onItemsPerPageChange={(newPageSize) => setPageSize(newPageSize)}
          onPageChange={(newPage) => setPage(newPage)}
        />
      ) : null}
    </div>
  );
};

export default ProcessList;
