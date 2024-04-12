// import { PlanType, Tiers } from './enums';

export const colors = {
  black: '#000000',
  blue01: '#d5e5fF',
  blue02: '#d5e5fe',
  blue03: '#538ade',
  blue04: '#343b40',
  blue05: '#021e40',
  blue06: '#1f60c4',
  blue07: '#245baf',
  blue08: '#001a3e',
  blue09: '#6898e2',
  blue10: '#3274d9',
  blue11: '#6e9fff',
  blue12: 'rgba(50, 116, 217, .15)',
  blue13: '#538CF1',
  blueBg01: '#1D2739',
  gray01: '#e9f4ff',
  gray02: '#345c97',
  gray03: '#cccfd2',
  gray04: '#7f7f82',
  green01: '#99d98d',
  green02: '#5aa64b',
  green03: '#1a7f4b',
  red01: '#ff7389',
  red02: '#de314d',
  orange01: '#ffb375',
  orange02: '#ff780a',
  orange03: '#ffb375',
  orange04: '#8A6C00',
  dark01: '#26262a',
  dark02: '#111217',
  dark03: '#424345',
  labelBorderColor01: 'rgba(204, 204, 220, 0.15)',
  labelTextColor01: '#ccccdc',
  labelBorderColor02: 'rgba(108, 207, 142, 0.5)',
  labelTextColor02: '#6ccf8e',
  labelBorderColor03: 'rgba(198, 155, 6, 0.5)',
  labelTextColor03: '#c69b06',
  labelBorderColor04: 'rgba(209, 14, 92, 0.5);',
  labelTextColor04: '#ff5286',
  labelBorderColor05: 'rgb(121, 158, 248)',
  labelTextColor05: '#799df8',

  alertColor: '#ff5286',
  yellow01: '#bb9205',
  filterBorderColor: '#ff8833',
};

export const labelColors = {
  dark: {
    labelBorderColor01: colors.labelBorderColor01,
    labelTextColor01: colors.labelTextColor01,
    labelBorderColor02: colors.labelBorderColor02,
    labelTextColor02: colors.labelTextColor02,
    labelBorderColor03: colors.labelBorderColor03,
    labelTextColor03: colors.labelTextColor03,
    labelBorderColor04: colors.labelBorderColor04,
    labelTextColor04: colors.labelTextColor04,
    labelBorderColor05: colors.labelBorderColor05,
    labelTextColor05: colors.labelTextColor05,
  },
  light: {
    labelBorderColor01: colors.gray04,
    labelTextColor01: colors.gray04,
    labelBorderColor02: colors.green02,
    labelTextColor02: colors.green02,
    labelBorderColor03: colors.yellow01,
    labelTextColor03: colors.yellow01,
    labelBorderColor04: colors.labelBorderColor04,
    labelTextColor04: colors.labelTextColor04,
    labelBorderColor05: colors.labelBorderColor05,
    labelTextColor05: colors.labelTextColor05,
  },
};

export const CHECK_TIMEOUT_MS = 60 * 1000;
export const CHECK_INTERVAL_MS = 5 * 1000;

export const AGENT_FAQ_LINK = 'https://grafana.com/docs/agent/latest/';

export const AGENT_CONFIG_NAME = 'agent-config.yaml';

// maybe this, too?
export const errorIntegrationApiResponse = {
  dashboard: 'unable to import dashboards: dashboard quota reached',
  alert: 'unable to apply alerting rules: could not create rule group: rule groups per tenant limit reached',
  record: 'unable to apply recording rules: could not create rule group: rule groups per tenant limit reached',
};

export const GRAFANA_EXAMPLE_USER = '<grafana.com username>';
export const GRAFANA_EXAMPLE_API = '<grafana.com API Key>';
export const DEFAULT_PROM_URL = 'https://prometheus-us-central1.grafana.net/api/prom/push';
export const DEFAULT_LOKI_URL = 'logs-prod-us-central1.grafana.net/api/prom/push';
export const DEFAULT_GRAPHITE_URL = 'https://graphite-us-central1.grafana.net/metrics';
export const DEFAULT_TEMPO_URL = 'tempo-us-central1.grafana.net';
export const DEFAULT_PYROSCOPE_URL = 'profiles-prod-001.grafana.net';
export const AWS_IAM_URL = 'https://console.aws.amazon.com/iam/home';

export const PLUGIN_ID = 'grafana-easystart-app';
export const PLUGIN_ID_STORAGE_KEY = `grafana.${PLUGIN_ID}.data`;

export const TERRAFORM_DOCS_URL =
  'https://grafana.com/docs/grafana-cloud/monitor-infrastructure/aws/cloudwatch-metrics/config-cw-metrics/#configure-with-terraform';

export const GRAFANA_AGENT_WINDOWS_URL =
  'https://github.com/grafana/agent/releases/latest/download/grafana-agent-installer.exe.zip';

export const GRAFANA_AGENT_WINDOWS_FLOW_URL =
  'https://github.com/grafana/agent/releases/latest/download/grafana-agent-flow-installer.exe.zip';

export enum FormErrors {
  REQUIRED_FIELD = 'This field is required',
  JOB_NAME_CHARACTERS = 'Scrape job name can only contain alphanumeric characters, dashes, and underscores.',
  SCRAPE_JOB_NAME_EXISTS = 'A scrape job with this name already exists. Please choose a unique name.',
}

export const archOptions = {
  Amd64: {
    label: 'Amd64',
    value: 'amd64',
  },
  Arm64: {
    label: 'Arm64',
    value: 'arm64',
  },
};

export const SCROLL_CONTAINER_SELECTOR = '.scrollbar-view';

export const DISABLED_REASONS = {
  disabled_by_user: {
    title: 'User disabled access',
    description: 'This scrape job was disabled by the user.',
  },
  credentials_revoked: {
    title: 'Credentials revoked',
    description: 'This scrape job was disabled because credentials have been deleted.',
  },
  over_series_limit: {
    title: 'Over series limit',
    description: 'This scrape job was disabled because the series limit has been reached.',
  },
};

export type DisabledReasonType = keyof typeof DISABLED_REASONS;

export const DEFAULT_SCRAPE_INTERVAL = '15s';
export const DEFAULT_SERVICES_SCRAPE_INTERVAL_SECONDS = 300;

export const CREATE_OR_UPDATE_JOBS_CACHE_KEY = 'create-or-update-jobs';

export const GRAFANA_AGENT_CHECK_ID = 'grafana-agent-check';
export const LINUX_NODE_ID = 'linux-node';
export const RASPBERRY_PI_NODE_ID = 'raspberry-pi-node';
export const CILIUM_ENTERPRISE_ID = 'cilium-enterprise';
export const AWS_ID = 'aws';
export const JAVA_ID = 'java';

export const grafanaAgentBaseUrl = 'https://github.com/grafana/agent/releases/latest/download';
export const grafanaAgentPrefix = 'grafana-agent';
export const DOCKER_DESKTOP_ID = 'docker-desktop';

export const ORG_INFO = 'grafanacloud_org_info';
export const PRO_PLAN_NAME = 'Pro plan';


export const HOSTNAME_RELABEL_KEY = '<your-instance-name>';

export const INSTALL_TROUBLESHOOTING_DOCS_LINK =
  'https://grafana.com/docs/grafana-cloud/monitor-infrastructure/integrations/install-troubleshoot';
