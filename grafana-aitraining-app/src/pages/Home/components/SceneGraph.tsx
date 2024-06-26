import React from 'react';

import { PanelData } from '@grafana/data';
import {
  EmbeddedScene,
  SceneFlexLayout,
  SceneFlexItem,
  PanelBuilders,
  SceneDataNode
//  SceneTimeRange,
} from '@grafana/scenes';
import { LineInterpolation, TooltipDisplayMode } from '@grafana/schema';

interface SceneGraphProps {
  data?: PanelData;
}

export const SceneGraph: React.FC<SceneGraphProps> = ({ data }) => {
  const myTimeSeriesPanel = PanelBuilders.timeseries().setTitle('My first panel');

  const dataNode = new SceneDataNode({ data });
  myTimeSeriesPanel.setData(dataNode);
  myTimeSeriesPanel.setOption('legend', { asTable: true }).setOption('tooltip', { mode: TooltipDisplayMode.Single });
  myTimeSeriesPanel.setDecimals(2).setUnit('ms');
  myTimeSeriesPanel.setCustomFieldConfig('lineInterpolation', LineInterpolation.Smooth);
  myTimeSeriesPanel.setOverrides((b) =>
    b.matchFieldsWithNameByRegex('/metrics/').overrideDecimals(4).overrideCustomFieldConfig('lineWidth', 5)
  );

  const myPanel = myTimeSeriesPanel.build();

  const scene = new EmbeddedScene({
    body: new SceneFlexLayout({
      children: [
        new SceneFlexItem({
          body: myPanel,
        }),
      ],
    }),
  });

  return <scene.Component model={scene} />;
};
