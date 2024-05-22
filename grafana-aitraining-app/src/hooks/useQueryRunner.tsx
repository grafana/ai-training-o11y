import React, { createContext, memo, useCallback, useContext, useEffect, useMemo } from 'react';
import { useObservable } from 'react-use';

import { PanelData, QueryRunner } from '@grafana/data';
import { createQueryRunner } from '@grafana/runtime';

const QueryContext = createContext<QueryRunner | undefined>(undefined);

type QueryRunnerProps = React.PropsWithChildren<{}> & { CustomContext?: React.Context<QueryRunner | undefined> };

export const QueryRunnerProvider = memo<QueryRunnerProps>(({ children, CustomContext }) => {
  const runner = useMemo(() => createQueryRunner(), []);

  useEffect(() => {
    const toDestroy = runner;
    return () => {
      return toDestroy.destroy();
    };
  }, [runner]);

  if (CustomContext !== undefined) {
    return <CustomContext.Provider value={runner}>{children}</CustomContext.Provider>;
  }

  return <QueryContext.Provider value={runner}>{children}</QueryContext.Provider>;
});

QueryRunnerProvider.displayName = 'QueryRunnerProvider';

export const useQueryRunner = (): QueryRunner => {
  const context = useContext(QueryContext);

  if (context == null) {
    throw new Error('You can only use `useQueryRunner` in a component wrapped in a `QueryRunnerProvider`.');
  }

  return context;
};

export const useQueryData = (): PanelData | undefined => {
  const runner = useQueryRunner();
  return useObservable(runner.get());
};

export const useCancelQuery = (): (() => void) => {
  const runner = useQueryRunner();

  return useCallback(() => {
    return runner.cancel();
  }, [runner]);
};
