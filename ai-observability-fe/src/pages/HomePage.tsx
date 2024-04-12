import React from 'react';
import { testIds } from '../components/testIds';
import { PluginPage } from '@grafana/runtime';
import {  Card, LinkButton, useStyles2 } from '@grafana/ui';

export function HomePage() {
  const styles = useStyles2(homeStyles);

  return (
    <PluginPage>
      <div data-testid={testIds.pageOne.container}>
        <div>
          <Card className={styles.cardCentered}>
            <h2>
              AI Observability is the multi-purpose, turnkey solution for monitoring your AI stack, streamlining triage, and
              optimizing resource utilization.
            </h2>
            <div>
              <LinkButton href='/plugins/gtm-aiobservability-app' style={{ marginRight: '10px' }} onClick={() => console.log('hi')}>
                Get started
              </LinkButton>
            </div>
          </Card>
          <Card className={styles.cardCentered}>
            <h2>Get insights on your:</h2>
            <div className={styles.cardGrid}>
              <div>
                <h3>LLMs</h3>
                <p>Monitor LLM token usage.</p>
              </div>
              <div>
                <h3>Vector Databases</h3>
                <p>
                  Monitor your vector databases, usage, and diagnostics.
                </p>
              </div>
              <div>
                <h3>ML Frameworks</h3>
                <p>
                  Monitor performance of your models.
                </p>
              </div>
              <div>
                <h3>Infrastructure</h3>
                <p>
                  Monitor performance of your GPUs
                </p>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </PluginPage>
  );
}


import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';

const homeStyles = (theme: GrafanaTheme2) => ({
  cardCentered: css`
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 40px;
    h2:first-of-type {
      display: none;
    }
    h2 {
      width: 60%;
      margin: 20px 0 40px 0;
    }
  `,
  cardGrid: css`
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    grid-template-rows: 1fr;
    width: 100%;
    text-align: center;

    h3 {
      margin-bottom: 20px;
    }

    img {
      margin-bottom: 20px;
    }

    p {
      // width: 300px;
      margin: auto;
    }
  `,
  card60: css`
    display: flex;
    flex-direction: column;
    width: 66.6%;
    padding: 40px;
  `,
  card30: css`
    width: 33.3%;
    margin-left: 8px;
    display: flex;
    flex-direction: column;
    padding: 40px;

    a {
      color: ${theme.colors.text.link};
    }

    ul {
      margin-left: 20px;
    }
  `,

});

export default homeStyles;