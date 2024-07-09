import React, { useEffect } from 'react';
import { useParams, useHistory } from 'react-router-dom';

import { GrafanaTheme2 } from '@grafana/data';
import { PluginPage } from '@grafana/runtime';
import { Tab, TabsBar, useStyles2 } from '@grafana/ui';

import { css } from '@emotion/css';

import { GraphsTab } from './components/GraphsTab';
import { ProcessList } from './components/ProcessList';
import { prefixRoute } from 'utils/utils.routing';
import { useTrainingAppStore } from 'utils/state';
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
    selectedRows,
    addSelectedRow,
    removeSelectedRow,
  } = trainingAppStore;
  const styles = useStyles2(getStyles);
  const params = useParams<{ path: string }>();
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
            // This block is actually unused right now because getBackendSrv in api.ts does not seem to propagate error codes
            switch (status) {
              case 401:
                setProcessesQueryStatus('unauthorized');
                break;
              case 404:
                setProcessesQueryStatus('notFound');
                break;
              case 500 || 502:
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
    <PluginPage
      renderTitle={() => {
        return (
          <div className={styles.pageHeader}>
            <div className={styles.pageTitle}>AI Training o11y</div>
          </div>
        );
      }}
    >
      <TabsBar>
        <Tab label="Process list" icon="table" active={tab === 'table'} onChangeTab={handleTabChange('table')} />
        <Tab
          label={selectedRows.length > 0 ? `Process graphs (${selectedRows.length})` : `Process graphs`}
          icon="graph-bar"
          active={tab === 'graphs'}
          onChangeTab={handleTabChange('graphs')}
        />
      </TabsBar>
      {tab === 'table' && (
        <ProcessList
          rows={renderedRows}
          processQueryStatus={processesQueryStatus}
          selectedRows={selectedRows}
          addSelectedRow={addSelectedRow}
          removeSelectedRow={removeSelectedRow}
        />
      )}
      {tab === 'graphs' && <GraphsTab rows={selectedRows} />}
    </PluginPage>
  );
};

const getStyles = (theme: GrafanaTheme2) => {
  return {
    tabsBar: css`
      margin-bottom: ${theme.spacing(3)};
    `,
    disabledTab: css`
      cursor: default;
      pointer-events: none;
      opacity: 0.2;
      background-color: red;
    `,
    pageHeader: css`
      display: flex;
      align-items: center;
      gap: 20px;
      height: 32px;
    `,
    pageTitle: css`
      font-family: Inter, Helvetica, Arial, sans-serif;
      font-size: 28px;
      font-weight: 500;
      line-height: 1.2;
      color: ${theme.colors.text.primary};
    `,
  };
};
