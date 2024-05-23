import React, { useEffect } from 'react';
import { Tab, TabsBar } from '@grafana/ui';
import { useParams } from 'react-router-dom';
import { PluginPage } from '@grafana/runtime';
import { prefixRoute } from 'utils/utils.routing';
import { PageLayoutType } from '@grafana/data';
import { useTrainingAppStore } from 'utils/state';
import { Graphs } from './components/Graphs';
import { TableTab } from './components/TableTab';
import { useGetProcesses } from 'utils/utils.plugin';

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

  const getProcesses = useGetProcesses();

  // This will fetch the processes from the backend
  // It will ultimately need to be more elaborate for paging, filtering, grouping, etc.
  useEffect(() => {
    const fetchProcesses = async () => {
      try {
        const response = await getProcesses();
        const data = response.data;
        setRenderedRows(data);
      } catch (error) {
        console.error('Error fetching processes:', error);
      }
    };

    fetchProcesses();
  }, [getProcesses, setRenderedRows]);

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
