import React from 'react';
import { testIds } from '../components/testIds';

import InstallDashboard from 'components/InstallDashboards/InstallDashboards';
import { grafanaAlloyPytorchServe, mlPytorchScrapeConfig } from './snippets';

export function MLFrameworks() {
  return (
    <div style={{ marginTop: '64px', border: '0.5px solid gray' }}>
      <h1 style={{ padding: '16px' }}> ML Framework Instructions</h1>
      <div data-testid={testIds.pageTwo.container}>
        <div style={{ marginLeft: '32px', padding: '16px' }}>
          <h1>Monitor PyTorch TorchServe</h1>
          <h2>Install Grafana Alloy</h2>
          <p>Grafana Alloy collects observability data and sends it to Grafana Cloud. Once the agent is deployed to your hosts, it collects and sends Prometheus-style metrics and log data using a pared-down Prometheus collector.</p>
          <p>Run this command to install and run Grafana Alloy as a grafana-alloy.service systemd service</p>
          <pre>{grafanaAlloyPytorchServe()}</pre>
          <h2>Prepare your configuration file</h2>
          <h3>Metrics:</h3>
          <p>Below <code>metrics.configs.scrape_configs</code>, insert the following lines and change the URLs according to your environment:</p>
          <pre>{mlPytorchScrapeConfig()}</pre>
          <h2>Restart the agent</h2>
          <p>Once you’ve made changes to your agent configuration file, run the following command to restart the agent.
          After installation, the Agent config is stored in /etc/grafana-agent.yaml. Restart the agent for any changes to take effect:</p>
          <pre>sudo systemctl restart grafana-agent.service</pre>

          <h2> Install dashboard</h2>
          <p>Get access to pre-configured dashboard that work right away</p>
        


          <InstallDashboard filePath={'https://raw.githubusercontent.com/grafana/hackathon-2024-03-tame-the-beast/main/gtm-aiobservability-app/src/instructions/ml/frameworks/pytorch.json?'} />
        </div>
      </div>
    </div>
  );
}
