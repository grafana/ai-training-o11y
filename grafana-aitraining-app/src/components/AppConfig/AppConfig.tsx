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
  stackId?: string;
  isMetadataTokenSet?: boolean;
};

type State = {
  metadataUrl: string;
  lokiDatasourceName: string;
  mimirDatasourceName: string;
  metadataToken: string;
  stackId: string;
  isMetadataTokenSet: boolean;
  dirty: {
    metadataUrl: boolean;
    lokiDatasourceName: boolean;
    mimirDatasourceName: boolean;
    metadataToken: boolean;
    stackId: boolean;
  };
};

interface Props extends PluginConfigPageProps<AppPluginMeta<JsonData>> {}

export const AppConfig = ({ plugin }: Props) => {
  const s = useStyles2(getStyles);
  const { enabled, pinned, jsonData } = plugin.meta;
  const originalValues = {
    metadataUrl: jsonData?.metadataUrl || '',
    lokiDatasourceName: jsonData?.lokiDatasourceName || '',
    mimirDatasourceName: jsonData?.mimirDatasourceName || '',
    stackId: jsonData?.stackId || '',
  };
  const [state, setState] = useState<State>({
    metadataUrl: jsonData?.metadataUrl || '',
    lokiDatasourceName: jsonData?.lokiDatasourceName || '',
    mimirDatasourceName: jsonData?.mimirDatasourceName || '',
    metadataToken: '',
    stackId: jsonData?.stackId || '',
    isMetadataTokenSet: plugin.meta.secureJsonFields?.metadataToken || false,
    dirty: {
      metadataUrl: false,
      lokiDatasourceName: false,
      mimirDatasourceName: false,
      metadataToken: false,
      stackId: false,
    },
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

  type StateKeys = keyof Omit<State, 'dirty'>;

  const onChangeField = (field: StateKeys, value: string | boolean) => {
    setState((prevState) => {
      // Only mark as dirty if the value has changed
      let isDirty = false;
    
      // It either is or is not set
      if (field !== 'isMetadataTokenSet') {
        isDirty = value !== originalValues[field as keyof typeof originalValues] && 
                  typeof value === 'string' && 
                  value.length > 0;
      }
      return {
        ...prevState,
        [field]: value,
        dirty: {
          ...prevState.dirty,
          [field]: isDirty,
        },
      };
    });
  };
  
  // Use this function for all field changes
  const onChangeMetadataToken = (event: ChangeEvent<HTMLInputElement>) => {
    onChangeField('metadataToken', event.target.value.trim());
  };
  
  const onChangeMetadataUrl = (event: ChangeEvent<HTMLInputElement>) => {
    onChangeField('metadataUrl', event.target.value.trim());
  };
  
  const onChangeLokiDatasource = (option: SelectableValue<string>) => {
    onChangeField('lokiDatasourceName', option.value || '');
  };
  
  const onChangeMimirDatasource = (option: SelectableValue<string>) => {
    onChangeField('mimirDatasourceName', option.value || '');
  };
  
  const onChangeStackId = (event: ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value.trim();
    if (/^\d*$/.test(value)) { // This regex ensures only digits are allowed
      onChangeField('stackId', value);
    }
  };

  // Function to check if the form should be enabled
const isFormEnabled = () => {
  return Object.entries(state.dirty).some(([field, isDirty]) => 
    field !== 'isMetadataTokenSet' && isDirty
  );
};

const handleSaveClick = () => {
  const isJsonDataDirty = ['metadataUrl', 'lokiDatasourceName', 'mimirDatasourceName', 'stackId'].some(
    field => state.dirty[field as keyof typeof state.dirty]
  );

  const isSecureJsonDataDirty = state.dirty.metadataToken;

  const updateData: any = {
    enabled,
    pinned,
  };

  if (isJsonDataDirty) {
    updateData.jsonData = {
      metadataUrl: state.metadataUrl,
      lokiDatasourceName: state.lokiDatasourceName,
      mimirDatasourceName: state.mimirDatasourceName,
      stackId: state.stackId,
    };
  }

  if (isSecureJsonDataDirty) {
    updateData.secureJsonData = {
      metadataToken: state.metadataToken,
    };
  }

  updatePluginAndReload(plugin.meta.id, updateData);
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
        <Field label="Metadata Service Token" description="A secret key for authenticating to our custom API">
          <SecretInput
            width={60}
            data-testid={testIds.appConfig.metadataToken}
            id="metadata-token"
            isConfigured={state.isMetadataTokenSet}
            value={state.metadataToken}
            placeholder={'Your secret API key'}
            onChange={onChangeMetadataToken}
            onReset={() => {
              onChangeField('metadataToken', '');
              onChangeField('isMetadataTokenSet', false);
            }}
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
            onClick={handleSaveClick}
            disabled={!isFormEnabled()}
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
