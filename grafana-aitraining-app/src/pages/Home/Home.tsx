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
    processesQueryStatus,
    setProcessesQueryStatus,
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
      setProcessesQueryStatus('loading');
      try {
        const response = await getProcesses();
        const data = response.data;
        setRenderedRows(data);
        setProcessesQueryStatus('success');
      } catch (error: unknown) {
        console.error('Error fetching processes:', error);
        
        if (error && typeof error === 'object') {
          if ('status' in error) {
            // Assuming the error object might have a status property
            const status = (error as { status: number }).status;
            switch (status) {
              case 401:
                setProcessesQueryStatus('unauthorized');
                break;
              case 404:
                setProcessesQueryStatus('notFound');
                break;
              case 500:
                setProcessesQueryStatus('serverError');
                break;
              default:
                setProcessesQueryStatus('error');
            }
          } else if ('message' in error) {
            // If there's a message property, log it
            console.error((error as { message: string }).message);
            setProcessesQueryStatus('error');
          } else {
            setProcessesQueryStatus('error');
          }
        } else {
          setProcessesQueryStatus('error');
        }
      }
    };
  
    fetchProcesses();
  }, [getProcesses, setRenderedRows, setProcessesQueryStatus]);

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
          processQueryStatus={processesQueryStatus}
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
