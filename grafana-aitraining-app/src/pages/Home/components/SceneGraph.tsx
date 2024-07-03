import React from 'react';

import { PanelData } from '@grafana/data';
import {
  SceneObjectState,
  SceneObjectBase,
  SceneComponentProps,
  EmbeddedScene,
  SceneFlexLayout,
  SceneFlexItem,
  VizPanel,
  PanelBuilders,
  SceneDataNode,
} from '@grafana/scenes';
import { LineInterpolation, TooltipDisplayMode } from '@grafana/schema';

export interface MetricPanel {
  pluginId: string;
  title: string;
  data: PanelData;
}
interface SceneGraphProps {
  panels?: MetricPanel[];
}

const getVizPanel = (data: PanelData, pluginId: string, title: string) => {
  const dataNode = new SceneDataNode({ data });
  return new SceneFlexItem({
    body: new VizPanel({
      pluginId,
      title,
      $data: dataNode,
    }),
    minHeight: 300,
    maxHeight: 300,
  });
};

export const SceneGraph: React.FC<SceneGraphProps> = ({ panels }) => {
  const scene = new EmbeddedScene({
    body: new SceneFlexLayout({
      children: panels?.map((p) => getVizPanel(p.data, p.pluginId, p.title)) || [],
    }),
  });
  return <scene.Component model={scene} />;
};
