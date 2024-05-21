import React, { useState } from 'react';
import { useAsync } from 'react-use';

import { dateTime, TimeRange } from '@grafana/data';
import { getDataSourceSrv } from '@grafana/runtime';
import { DataQuery } from '@grafana/schema';

import { useQueryResult } from 'hooks/useQueryResults';

interface LokiProps {
  datasource: any;
  data?: any;
  timeRange: TimeRange;
  onRunQuery: (query: DataQuery) => void;
  onQueryResult: (data: any) => void;
}

type LokiWrapperProps = Pick<LokiProps, 'onQueryResult'>;

const LokiWrapper: React.FC<LokiWrapperProps> = ({ onQueryResult }) => {

  // get here grabs the datasource with the component attached
  const lokiDS = useAsync(async () => {
    return getDataSourceSrv().get('Loki');
  }, [name]);

  // need to set actual real timeranges in ui
  const startDate = new Date();
  const endDate = dateTime(new Date());
  const tmpTimeRange: TimeRange = {
    from: dateTime(startDate),
    to: endDate,
    raw: { from: startDate.toLocaleString(), to: endDate.toLocaleString() },
  };

  // data, runQuery, cancelQuery
  const [data,onRunQuery,,isRunning] = useQueryResult(
    100,
    tmpTimeRange,
    'EST',
    lokiDS
  );

  // TODO: Remove this - Example of using isRunning to determine when the results of the query are ready
  console.log('isRunning', isRunning);

  if ((lokiDS?.value ?? '') === '') {
    console.log('Error loading datasource editor: NULL DATASOURCE');
    return <div>Error loading datasource editor</div>;
  }

  return (
    <div>
      <Loki
        datasource={lokiDS}
        data={data}
        timeRange={tmpTimeRange}
        onRunQuery={onRunQuery}
        onQueryResult={onQueryResult}
      />
    </div>
  );
};

const Loki: React.FC<LokiProps> = ({ datasource, data, timeRange, onRunQuery, onQueryResult }) => {
  const [query, setQuery] = useState<DataQuery>({ refId: 'A' });

  const ReactQueryEditor = datasource.value?.components?.QueryEditor;

  if (ReactQueryEditor == null) {
    return <div>Query editor not available for datasource</div>;
  }

  return (
    <ReactQueryEditor
      data={data}
      range={timeRange}
      onRunQuery={() => {
        console.log('running query', query);
        onRunQuery(query);
      }}
      query={query}
      datasource={datasource.value!}
      onChange={(q: DataQuery) => setQuery(q)}
    />
  );
};

export { LokiWrapper as Loki };
