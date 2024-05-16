import React, { useEffect } from 'react';
import { Tab, TabsBar } from '@grafana/ui';
import { useParams } from 'react-router-dom';
import { PluginPage } from '@grafana/runtime';
import { prefixRoute } from 'utils/utils.routing';
import { PageLayoutType } from '@grafana/data';
import { useTabStore } from 'utils/state';
import { Graphs } from './components/Graphs';
import { Table } from './components/Table';

// Placeholders for development, if this is in a PR deny it
const tableData = [
  { column1: 'Row 1, Column 1', column2: 'Row 1, Column 2' },
  { column1: 'Row 2, Column 1', column2: 'Row 2, Column 2' },
  { column1: 'Row 3, Column 1', column2: 'Row 3, Column 2' },
];

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

  useEffect(() => {
    setTab(tabFromUrl);
  }, [tabFromUrl, setTab]);


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
      {getTab === 'table' && <Table data={tableData} />}
      {getTab === 'graphs' && <Graphs data={graphsData} />}
    </PluginPage>
  );
};
