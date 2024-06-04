import { getDataSourceSrv } from '@grafana/runtime';
import { /* DataQueryRequest, */ TimeRange, dateTime } from '@grafana/data';
import React from 'react';
import { useAsync } from 'react-use';
import { RowData } from 'utils/state';
// import { DataQuery } from '@grafana/schema';

interface GraphsProps {
  rows: RowData[];
}

export const GraphsTab: React.FC<GraphsProps> = ({ rows }) => {

  const lokiDS = useAsync(async () => {
    return getDataSourceSrv().get('Loki');
  }, [name]);

  const { loading, error, value: dataSource } = lokiDS;

  if (loading || error || !dataSource) {
    // Handle loading state or error
    return null;
  }

  const startDate = new Date();
  startDate.setHours(startDate.getHours() - 7200);
  const endDate = dateTime(new Date());
  const tmpTimeRange: TimeRange = {
    from: dateTime(startDate),
    to: endDate,
    raw: { from: startDate.toLocaleString(), to: endDate.toLocaleString() },
  };

  console.log(tmpTimeRange);
  
  // const request: DataQueryRequest<DataQuery> = {
  //   targets: [
  //     {
  //       refId: 'A',
  //       queryType: 'logql',
  //       query: '{job="o11y"}',
  //       // Add other necessary properties for the Loki query
  //     },
  //   ],
  //   range: tmpTimeRange,
  //   // Add other properties as needed
  // };
  
  // dataSource.query(request).then((response: any) => {
  //   // Handle the query response
  //   console.log(response);
  // });

  return (
    <div>
      <h2>Graphs Props</h2>
      <pre>{JSON.stringify(rows, null, 2)}</pre>
    </div>
  );
};
