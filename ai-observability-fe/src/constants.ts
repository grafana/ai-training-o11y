import pluginJson from './plugin.json';

export const PLUGIN_BASE_URL = `/a/${pluginJson.id}`;

export enum ROUTES {
  HomePage = 'home',
  MLFrameworks = 'ml',
  Infrastructure = 'infra',
  VectorDB = 'vectordb',
  LLM = 'llm',
  CONFIG = 'config',
}
