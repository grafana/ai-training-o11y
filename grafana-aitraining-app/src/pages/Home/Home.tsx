import React, { useEffect } from 'react';
import { Tab, TabsBar } from '@grafana/ui';
import { useParams } from 'react-router-dom';
import { PluginPage } from '@grafana/runtime';
import { prefixRoute } from 'utils/utils.routing';
import { PageLayoutType } from '@grafana/data';
import { RowData, useTrainingAppStore } from 'utils/state';
import { Graphs } from './components/Graphs';
import { TableTab } from './components/TableTab';

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

export const Home = () => {
  // state management
  const trainingAppStore = useTrainingAppStore();
  const {
    tab,
    setTab,
    renderedRows,
    setRenderedRows,
    isSelected,
    setIsSelected,
    selectedRows,
    addSelectedRow,
    removeSelectedRow,
  } = trainingAppStore;

  const params = useParams<{path: string}>();
  // What tab we are on
  let tabFromUrl = params['path']?.split('/')[0];
  tabFromUrl = tabFromUrl === 'table' || tabFromUrl === 'graphs' ? tabFromUrl : 'table';
  useEffect(() => {
    setTab(tabFromUrl as "table" | "graphs");
  }, [tabFromUrl, setTab]);

  // This will need to be made more elaborate for paging, filtering, grouping, etc.
  useEffect(() => {
    // Fetch table data and update the rows state
    const tableData = fetchTableData();
    setRenderedRows(tableData);
  }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  , []);

  return (
    <PluginPage
      layout={PageLayoutType.Canvas}
    >
      <TabsBar>
        <Tab label="Process table" icon='table' href={prefixRoute('table')} active={tab === 'table'}>
        </Tab>
        <Tab label="Process graphs" icon='graph-bar' href={prefixRoute('graphs')} active={tab === 'graphs'}>
        </Tab>
      </TabsBar>
      {tab === 'table' &&
        <TableTab
          rows={renderedRows}
          isSelected={isSelected}
          setIsSelected={setIsSelected}
          addSelectedRow={addSelectedRow}
          removeSelectedRow={removeSelectedRow}
        />
      }
      {tab === 'graphs' && <Graphs rows={selectedRows} />}
    </PluginPage>
  );
};
