import React, { useEffect } from 'react';
import { Tab, TabsBar } from '@grafana/ui';
import { useParams } from 'react-router-dom';
import { PluginPage } from '@grafana/runtime';
import { prefixRoute } from 'utils/utils.routing';
import { PageLayoutType } from '@grafana/data';
import { useTabStore, useRowsStore, RowData } from 'utils/state';
import { Graphs } from './components/Graphs';
import { Table } from './components/Table';

// Function to fetch table data (placeholder for development)
const fetchTableData: () => RowData[] = () => {
  return [
    {
      process_uuid: 'abc123',
      status: 'running',
      start_time: '2023-06-08T10:00:00',
      additional_field1: 'Value 1',
    },
    {
      process_uuid: 'def456',
      status: 'finished',
      start_time: '2023-06-08T11:30:00',
      end_time: '2023-06-08T12:15:00',
      additional_field2: 'Value 2',
    },
    {
      process_uuid: 'ghi789',
      status: 'crashed',
      start_time: '2023-06-08T13:45:00',
      end_time: '2023-06-08T14:00:00',
      additional_field3: 'Value 3',
    },
    {
      process_uuid: 'jkl012',
      status: 'timed out',
      start_time: '2023-06-08T15:20:00',
      end_time: '2023-06-08T16:00:00',
      additional_field4: 'Value 4',
    },
  ];
};


const graphsData = [
  { value: 10 },
  { value: 20 },
  { value: 15 },
  { value: 25 },
  { value: 30 },
];

export const Home = () => {
  const params = useParams<{path: string}>();

  // What tab we are on
  let tabFromUrl = params['path']?.split('/')[0];
  tabFromUrl = tabFromUrl === 'table' || tabFromUrl === 'graphs' ? tabFromUrl : 'table';
  const setTab = useTabStore((state) => state.set);
  const getTab = useTabStore((state) => state.tab);

  // Access the rows state and setRows function from useRowsStore
  const rows = useRowsStore((state) => state.rows);
  const setRows = useRowsStore((state) => state.setRows);

  useEffect(() => {
    setTab(tabFromUrl);
  }, [tabFromUrl, setTab]);

  useEffect(() => {
    if (getTab === 'table') {
      // Fetch table data and update the rows state
      const tableData = fetchTableData();
      setRows(tableData);
    }
  }, [getTab, setRows]);

  return (
    <PluginPage
      layout={PageLayoutType.Canvas}
    >
      <TabsBar>
        <Tab label="Process table" icon='table' href={prefixRoute('table')} active={getTab === 'table'}>
        </Tab>
        <Tab label="Process graphs" icon='graph-bar' href={prefixRoute('graphs')} active={getTab === 'graphs'}>
        </Tab>
      </TabsBar>
      {getTab === 'table' && <Table data={rows} />}
      {getTab === 'graphs' && <Graphs data={graphsData} />}
    </PluginPage>
  );
};
