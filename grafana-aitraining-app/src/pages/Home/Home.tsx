import React, { useEffect } from 'react';
import { Tab, TabsBar } from '@grafana/ui';
import { useParams, useHistory } from 'react-router-dom';
import { PluginPage } from '@grafana/runtime';
import { prefixRoute } from 'utils/utils.routing';
import { PageLayoutType } from '@grafana/data';
import { useTrainingAppStore } from 'utils/state';
import { GraphsTab } from './components/GraphsTab';
import { TableTab } from './components/TableTab';
import { useGetProcesses } from 'utils/utils.plugin';

export const Home = () => {
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
  const history = useHistory();
  
  let tabFromUrl = params['path'] as 'table' | 'graphs';

  useEffect(() => {
    if (tabFromUrl === 'graphs' && (!selectedRows || selectedRows.length === 0)) {
      history.replace(prefixRoute('table'));
    } else {
      setTab(tabFromUrl);
    }
  }, [tabFromUrl, selectedRows, history, setTab]);

  const getProcesses = useGetProcesses();

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

  const handleTabChange = (newTab: 'table' | 'graphs') => (event: React.MouseEvent<HTMLElement>) => {
    event.preventDefault();
    if (newTab === 'graphs' && (!selectedRows || selectedRows.length === 0)) {
      // Do nothing or show a message that rows need to be selected first
      return;
    }
    history.push(prefixRoute(newTab));
  };

  return (
    <PluginPage layout={PageLayoutType.Canvas}>
      <TabsBar>
        <Tab 
          label="Process table" 
          icon="table"
          active={tab === 'table'}
          onChangeTab={handleTabChange('table')}
        />
        <Tab 
          label="Process graphs" 
          icon="graph-bar"
          active={tab === 'graphs'}
          onChangeTab={handleTabChange('graphs')}
        />
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
      {tab === 'graphs' && <GraphsTab rows={selectedRows} />}
    </PluginPage>
  );
};
