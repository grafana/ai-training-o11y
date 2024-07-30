import React, { useState, useEffect, ChangeEvent } from 'react';
import { Button, Field, Input, useStyles2, FieldSet, SecretInput, Select } from '@grafana/ui';
import { PluginConfigPageProps, AppPluginMeta, PluginMeta, GrafanaTheme2, SelectableValue } from '@grafana/data';
import { getBackendSrv, locationService } from '@grafana/runtime';
import { css } from '@emotion/css';
import { testIds } from '../testIds';
import { lastValueFrom } from 'rxjs';

export type JsonData = {
  metadataUrl?: string;
  lokiDatasourceName?: string;
  mimirDatasourceName?: string;
  isApiKeySet?: boolean;
  stackId?: string;
};

type State = {
  metadataUrl: string;
  lokiDatasourceName: string;
  mimirDatasourceName: string;
  isApiKeySet: boolean;
  apiKey: string;
  stackId: string;
};

interface Props extends PluginConfigPageProps<AppPluginMeta<JsonData>> {}

export const AppConfig = ({ plugin }: Props) => {
  const s = useStyles2(getStyles);
  const { enabled, pinned, jsonData } = plugin.meta;
  const [state, setState] = useState<State>({
    metadataUrl: jsonData?.metadataUrl || '',
    lokiDatasourceName: jsonData?.lokiDatasourceName || '',
    mimirDatasourceName: jsonData?.mimirDatasourceName || '',
    apiKey: '',
    isApiKeySet: Boolean(jsonData?.isApiKeySet),
    stackId: jsonData?.stackId || '',
  });

  // eslint-disable-next-line @typescript-eslint/array-type
  const [datasources, setDatasources] = useState<Array<SelectableValue<string>>>([]);

  useEffect(() => {
    const fetchDatasources = async () => {
      try {
        const sources = await getDatasources();
        setDatasources(sources);
      } catch (error) {
        console.error('Error fetching datasources:', error);
      }
    };

    fetchDatasources();
  }, []);

  const onResetApiKey = () =>
    setState({
      ...state,
      apiKey: '',
      isApiKeySet: false,
    });

  const onChangeApiKey = (event: ChangeEvent<HTMLInputElement>) => {
    setState({
      ...state,
      apiKey: event.target.value.trim(),
    });
  };

  const onChangeMetadataUrl = (event: ChangeEvent<HTMLInputElement>) => {
    setState({
      ...state,
      metadataUrl: event.target.value.trim(),
    });
  };

  const onChangeLokiDatasource = (option: SelectableValue<string>) => {
    setState({
      ...state,
      lokiDatasourceName: option.value || '',
    });
  };

  const onChangeMimirDatasource = (option: SelectableValue<string>) => {
    setState({
      ...state,
      mimirDatasourceName: option.value || '',
    });
  };

  const onChangeStackId = (event: ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value.trim();
    if (/^\d*$/.test(value)) { // This regex ensures only digits are allowed
      setState({
        ...state,
        stackId: value,
      });
    }
  };

  return (
    <div data-testid={testIds.appConfig.container}>
      {/* ENABLE / DISABLE PLUGIN */}
      <FieldSet label="Enable / Disable">
        {!enabled && (
            <>
              <div className={s.colorWeak}>The plugin is currently not enabled.</div>
              <Button
                className={s.marginTop}
                variant="primary"
                onClick={() =>
                  updatePluginAndReload(plugin.meta.id, {
                    enabled: true,
                    pinned: true,
                    jsonData,
                  })
                }
              >
                Enable plugin
              </Button>
            </>
          )}

          {/* Disable the plugin */}
          {enabled && (
            <>
              <div className={s.colorWeak}>The plugin is currently enabled.</div>
              <Button
                className={s.marginTop}
                variant="destructive"
                onClick={() =>
                  updatePluginAndReload(plugin.meta.id, {
                    enabled: false,
                    pinned: false,
                    jsonData,
                  })
                }
              >
                Disable plugin
              </Button>
            </>
          )}
        </FieldSet>

      {/* CUSTOM SETTINGS */}
      <FieldSet label="API Settings" className={s.marginTopXl}>
        {/* API Key */}
        <Field label="API Key" description="A secret key for authenticating to our custom API">
          <SecretInput
            width={60}
            data-testid={testIds.appConfig.apiKey}
            id="api-key"
            value={state.apiKey}
            isConfigured={state.isApiKeySet}
            placeholder={'Your secret API key'}
            onChange={onChangeApiKey}
            onReset={onResetApiKey}
          />
        </Field>

        {/* Metadata URL */}
        <Field label="Metadata URL" description="URL for the metadata API" className={s.marginTop}>
          <Input
            width={60}
            id="metadata-url"
            data-testid={testIds.appConfig.metadataUrl}
            value={state.metadataUrl}
            placeholder="http://ai-training-api:8000"
            onChange={onChangeMetadataUrl}
          />
        </Field>

        {/* Loki Datasource */}
        <Field label="Loki Datasource" description="Select the Loki datasource" className={s.marginTop}>
          <Select
            width={60}
            id="loki-datasource"
            data-testid={testIds.appConfig.lokiDatasource}
            value={state.lokiDatasourceName}
            onChange={onChangeLokiDatasource}
            options={datasources}
            placeholder="Select Loki datasource"
          />
        </Field>

        {/* Mimir Datasource */}
        <Field label="Mimir Datasource" description="Select the Mimir datasource" className={s.marginTop}>
          <Select
            width={60}
            id="mimir-datasource"
            data-testid={testIds.appConfig.mimirDatasource}
            value={state.mimirDatasourceName}
            onChange={onChangeMimirDatasource}
            options={datasources}
            placeholder="Select Mimir datasource"
          />
        </Field>

        {/* Stack ID */}
        <Field label="Stack ID" description="Enter the numeric stack ID" className={s.marginTop}>
          <Input
            width={60}
            id="stack-id"
            data-testid={testIds.appConfig.stackId}
            value={state.stackId}
            placeholder="Enter stack ID (only digits allowed)"
            onChange={onChangeStackId}
          />
        </Field>

        <div className={s.marginTop}>
          <Button
            type="submit"
            data-testid={testIds.appConfig.submit}
            onClick={() =>
              updatePluginAndReload(plugin.meta.id, {
                enabled,
                pinned,
                jsonData: {
                  metadataUrl: state.metadataUrl,
                  lokiDatasourceName: state.lokiDatasourceName,
                  mimirDatasourceName: state.mimirDatasourceName,
                  isApiKeySet: true,
                  stackId: state.stackId,
                },
                secureJsonData: state.isApiKeySet
                  ? undefined
                  : {
                      apiKey: state.apiKey,
                    },
              })
            }
            disabled={Boolean(
              !state.metadataUrl ||
                !state.lokiDatasourceName ||
                !state.mimirDatasourceName ||
                (!state.isApiKeySet && !state.apiKey) ||
                !state.stackId
            )}
          >
            Save API settings
          </Button>
        </div>
      </FieldSet>
    </div>
  );
};

// Implement this function to fetch available datasources
// eslint-disable-next-line @typescript-eslint/array-type
const getDatasources = async (): Promise<Array<SelectableValue<string>>> => {
  try {
    const response = await getBackendSrv().get('/api/datasources');
    return response.map((ds: any) => ({
      label: ds.name,
      value: ds.name,
    }));
  } catch (error) {
    console.error('Error fetching datasources:', error);
    return [];
  }
};

const getStyles = (theme: GrafanaTheme2) => ({
  colorWeak: css`
    color: ${theme.colors.text.secondary};
  `,
  marginTop: css`
    margin-top: ${theme.spacing(3)};
  `,
  marginTopXl: css`
    margin-top: ${theme.spacing(6)};
  `,
});

const updatePluginAndReload = async (pluginId: string, data: Partial<PluginMeta<JsonData>>) => {
  try {
    await updatePlugin(pluginId, data);

    // Reloading the page as the changes made here wouldn't be propagated to the actual plugin otherwise.
    // This is not ideal, however unfortunately currently there is no supported way for updating the plugin state.
    locationService.reload();
  } catch (e) {
    console.error('Error while updating the plugin', e);
  }
};

export const updatePlugin = async (pluginId: string, data: Partial<PluginMeta>) => {
  const response = getBackendSrv().fetch({
    url: `/api/plugins/${pluginId}/settings`,
    method: 'POST',
    data,
  });

  const dataResponse = await lastValueFrom(response);

  return dataResponse.data;
};
