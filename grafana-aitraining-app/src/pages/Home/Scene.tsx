import { EmbeddedScene, PanelBuilders, SceneApp, SceneAppPage, SceneDataTransformer, SceneFlexItem, SceneFlexLayout, SceneQueryRunner } from '@grafana/scenes';
import { ROUTES } from '../../constants';
import { prefixRoute } from '../../utils/utils.routing';

export const MY_DATASOURCE_REF = {
  uid: 'grafana-aitraining-app-datasource-uid'
};

export function getProcessesTable() {
  const data = new SceneDataTransformer({
    transformations: []});

  return PanelBuilders.table()
    .setTitle('Processes')
    .setData(data)
    .setHoverHeader(true)
    .setDisplayMode('transparent')
    .setOption('sortBy', [{ displayName: 'start_time' }])
    .build();
}


const getTab1Scene = () =>
  new EmbeddedScene({ 
    $data: new SceneQueryRunner({
      datasource: MY_DATASOURCE_REF,
      queries: [
        { refId: 'A', datasource: MY_DATASOURCE_REF}
      ],
    }),

    body: new SceneFlexLayout({
      direction: 'column',
      children: [
        new SceneFlexItem({
          height: 300,
          body: getProcessesTable(),
        }),
      ],
    }),
  });

const getTab2Scene = () => {
  return new EmbeddedScene({
    body: new SceneFlexLayout({
      children: [
        new SceneFlexItem({
          width: '100%',
          height: 300,
          body: PanelBuilders.text().setTitle('Hello world panel').setOption('content', 'Hello world!').build(),
        }),
      ],
    }),
  });};

export const getScene = () => {
  return new SceneApp({
    pages: [
      new SceneAppPage({
        title: 'AI Training Observability',
        // Important: Mind the page route is unambiguous for the tabs to work properly
        url: prefixRoute(`${ROUTES.Home}`),
        hideFromBreadcrumbs: true,
        getScene: getTab1Scene,
        tabs: [
          new SceneAppPage({
            title: 'Processes',
            url: prefixRoute(`${ROUTES.Home}`),
            getScene: getTab1Scene,
          }),
          new SceneAppPage({
            title: 'Graphs',
            url: prefixRoute(`${ROUTES.Home}/graphs`),
            getScene: getTab2Scene,
          }),
        ],
      }),
    ],
  });
}
