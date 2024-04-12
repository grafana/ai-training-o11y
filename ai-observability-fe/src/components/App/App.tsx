import React from 'react';
import { Route, Routes } from 'react-router-dom';
import { AppRootProps } from '@grafana/data';
import { ROUTES } from '../../constants';
import { VectorDB, HomePage, Infrastructure, MLFrameworks, LLM } from '../../pages';

export function App(props: AppRootProps) {
  return (
    <Routes>
      <Route path={ROUTES.MLFrameworks} element={<MLFrameworks />} />
      <Route path={`${ROUTES.Infrastructure}/:id?`} element={<Infrastructure />} />
      <Route path={ROUTES.LLM} element={<LLM />} />


      {/* Full-width page (this page will have no side navigation) */}
      <Route path={ROUTES.VectorDB} element={<VectorDB />} />

      {/* Default page */}
      <Route path="*" element={<HomePage />} />
    </Routes>
  );
}
