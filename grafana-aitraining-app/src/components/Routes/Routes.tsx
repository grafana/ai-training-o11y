import React from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';
import { Home } from '../../pages/Home';
import { prefixRoute } from '../../utils/utils.routing';
import { ROUTES } from '../../constants';

export const Routes = () => {
  return (
    <Switch>
      <Route path={prefixRoute(`${ROUTES.Home}`)} component={Home} />
      <Redirect to={prefixRoute(ROUTES.Home)} />
    </Switch>
  );
};
