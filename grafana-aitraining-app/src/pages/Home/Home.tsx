import React from 'react';
import { LinkButton, useStyles2 } from '@grafana/ui';
import { useParams } from 'react-router-dom';
import { PluginPage } from '@grafana/runtime';
import { testIds } from '../../components/testIds';
import { prefixRoute } from 'utils/utils.routing';
import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';


export const Home = () => {
  const params = useParams();
  console.log('Path arguments:', params);
  const s = useStyles2(getStyles);


  return (
    <PluginPage>
      <div data-testid={testIds.pageOne.container}>
        This is page one.
        <div className={s.marginTop}>
          <LinkButton data-testid={testIds.pageOne.navigateToFour} href={prefixRoute('two')}>
            Full-width page example
          </LinkButton>
        </div>
      </div>
    </PluginPage>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  marginTop: css`
    margin-top: ${theme.spacing(2)};
  `,
});
