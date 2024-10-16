import React, { useEffect } from 'react';
import { useParams, useHistory } from 'react-router-dom';

import { GrafanaTheme2 } from '@grafana/data';
import { PluginPage } from '@grafana/runtime';
import { useStyles2, Button, TextLink, Alert, Spinner } from '@grafana/ui';

import { css } from '@emotion/css';

import { GraphsView } from './components/GraphsView';
import { ProcessList } from './components/ProcessList';
import { prefixRoute } from 'utils/utils.routing';
import { useTrainingAppStore } from 'utils/state';
import { useGetProcesses } from 'utils/utils.plugin';
import { useSettings } from 'hooks/useSettings';

export const Home = () => {
  const trainingAppStore = useTrainingAppStore();
  const {
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
  const {isReady, error, settings } = useSettings()

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

  const handleViewChange = (view: 'table' | 'graphs') => {
    history.push(prefixRoute(view));
  };

  if (!isReady) {
    return <Spinner />;
  }

  if (error !== undefined || settings === undefined) {
    return <Alert title="Error loading"> Error loading app settings. Please try again later. </Alert>;
  }

  const isCloud = settings.stackId !== undefined && settings.stackId !== '';
  let metadataUrl = new URL(settings.metadataUrl);
  if (isCloud) {
    metadataUrl.username = settings.stackId
    // Setting the password appears to pass through URL encoding still so <token> is not rendered properly.
    metadataUrl.password = '--token--'
  }

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
      {tabFromUrl === 'table' ? (
        <>
          <div className={styles.instructionBox}>
            <div>
              Select one or more training processes from the list below.
              <br />
              Then, click <strong>View graphs</strong> and a list of graphs will be generated.
              <p>
                To send data to the AI Training o11y app set the following environment variable:
              </p>
              <pre>
                {'GF_AI_TRAINING_CREDS="' + metadataUrl + '"'}
              </pre>
              {isCloud ? (
                // Only generate this section if a stackId is present (running in cloud).
                <p>
                  To generate a token create an access policy in the <TextLink href="/a/grafana-auth-app">Grafana auth app</TextLink>.
                </p>
              ) : null}
            </div>
            {selectedRows.length < 1 ? (
              <Button disabled={true} variant="primary">
                View graphs (0 processes)
              </Button>
            ) : (
              <Button
                onClick={() => {
                  console.log('hey');
                  handleViewChange('graphs');
                }}
                variant="primary"
              >
                View graphs ({selectedRows.length} processes)
              </Button>
            )}
          </div>

          <ProcessList
            rows={renderedRows}
            processQueryStatus={processesQueryStatus}
            selectedRows={selectedRows}
            addSelectedRow={addSelectedRow}
            removeSelectedRow={removeSelectedRow}
          />
        </>
      ) : (
        <div>
          <div className={styles.instructionBox}>
          <div>
              Review graphs below for the processes you selected
            </div>            
            <Button onClick={() => handleViewChange('table')}>
              Back to process list
            </Button>
          </div>
          <GraphsView rows={selectedRows} />
        </div>
      )}
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
    instructionBox: css`
      display: flex;
      justify-content: space-between;
      padding: 20px;
      background-color: rgba(61, 113, 217, 0.15);
      border: solid 1px rgba(110, 159, 255, 0.25);
    `,
  };
};
