import React from 'react';
import { getScene } from './Scene';

export const Home = () => {
  const scene = getScene();

  return <scene.Component model={scene} />;
};
