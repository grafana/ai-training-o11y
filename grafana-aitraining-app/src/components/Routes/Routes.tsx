import React from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';
import { Home } from '../../pages/Home';
import { prefixRoute } from '../../utils/utils.routing';

export const Routes = () => {
  return (
    <Switch>
      <Route path={prefixRoute(':path(table|graphs)')} component={Home} />
      <Redirect to={prefixRoute('table')} />
    </Switch>
  );
};
