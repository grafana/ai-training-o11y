import React, { useEffect, useRef } from 'react';

import { useProcessQueries } from 'hooks/useProcessQueries';
import { useTrainingAppStore, RowData } from 'utils/state';
import { reshapeModelMetrics } from 'utils/reshapeModelMetrics';
import { SceneGraph } from './SceneGraph';

import { PanelData, DataFrame, LoadingState, dateTime, TimeRange, FieldType } from '@grafana/data';

export interface MetricPanel {
  pluginId: string;
  title: string;
  data: PanelData;
}
interface GraphsProps {
  rows: RowData[];
}

export const GraphsTab: React.FC<GraphsProps> = ({ rows }) => {
  const {
    lokiQueryStatus,
    lokiQueryData,
    organizedLokiData,
    resetLokiResults,
    setOrganizedLokiData
  } = useTrainingAppStore();
  const { isReady, runQueries } = useProcessQueries();
  const shouldRunQueries = useRef(true);

  useEffect(() => {
    if (isReady && rows.length > 0 && shouldRunQueries.current) {
      shouldRunQueries.current = false;
      resetLokiResults();
      runQueries();
    }
  }, [isReady, rows, resetLokiResults, runQueries]);

  useEffect(() => {
    shouldRunQueries.current = true;
  }, [rows]);

  useEffect(() => {
    if (lokiQueryStatus === 'success' && Object.keys(lokiQueryData).length > 0) {
      const organized = reshapeModelMetrics(lokiQueryData);
      setOrganizedLokiData(organized);
    }
  }, [lokiQueryStatus, lokiQueryData, setOrganizedLokiData]);

  if (!isReady) {
    return <div>Loading...</div>;
  }

  if (lokiQueryStatus === 'loading') {
    return (
      <div>
        Running...
        <button onClick={() => { resetLokiResults(); shouldRunQueries.current = true; }}>Reset Results</button>
      </div>
    );
  }

  if (!organizedLokiData) {
    return <div>No data</div>;
  }

  /// ---- FAKE DATA BELOW -----

  /// FAKE TIME RANGE
  const tmpTimeRange: TimeRange = {
    from: dateTime('2024-06-26T00:01:00.001Z'), // startDate,
    to: dateTime('2024-06-26T10:30:00.001Z'), // endDate,
    raw: {
      from: dateTime('2024-06-26T00:01:00.001Z').toISOString(),
      to: dateTime('2024-06-26T10:30:00.001Z').toISOString(),
    },
  };

  /// FAKE DATA FRAME
  const sampleFrame: DataFrame = {
    name: '',
    fields: [
      {
        name: 'TimeType',
        type: FieldType.time,
        values: [
          1719367200000, 1719374400000, 1719381600000, 1719388800000, 1719396000000, 1719403200000, 1719410400000,
          1719417600000, 1719424800000, 1719432000000, 1719439200000, 1719446400000, 1719453600000, 1719460800000,
          1719468000000, 1719475200000, 1719482400000, 1719489600000, 1719496800000, 1719504000000, 1719511200000,
          1719518400000, 1719525600000, 1719532800000, 1719540000000, 1719547200000, 1719554400000, 1719561600000,
          1719568800000, 1719576000000, 1719583200000, 1719590400000, 1719597600000, 1719604800000, 1719612000000,
          1719619200000, 1719626400000, 1719633600000, 1719640800000, 1719648000000, 1719655200000, 1719662400000,
          1719669600000, 1719676800000, 1719684000000, 1719691200000, 1719698400000, 1719705600000, 1719712800000,
          1719720000000, 1719727200000, 1719734400000, 1719741600000, 1719748800000, 1719756000000, 1719763200000,
          1719770400000, 1719777600000, 1719784800000, 1719792000000, 1719799200000, 1719806400000, 1719813600000,
          1719820800000, 1719828000000, 1719835200000, 1719842400000, 1719849600000, 1719856800000, 1719864000000,
          1719871200000, 1719878400000, 1719885600000, 1719892800000, 1719900000000, 1719907200000, 1719914400000,
          1719921600000, 1719928800000, 1719936000000, 1719943200000, 1719950400000, 1719957600000, 1719964800000,
          1719972000000,
        ],
        config: {},
      },
      { name: 'NumberType1', type: FieldType.number, values: [1, 2, 3, 4, 5, 6, 7, 8], config: {} },
      { name: 'NumberType2', type: FieldType.number, values: [4, 4, 4, 4, 4, 2, 4, 2], config: {} },
    ],
    length: 2,
  };

  /// FAKE DATA PANE
  const samplePanel: PanelData = {
    state: LoadingState.Done,
    series: [sampleFrame],
    timeRange: tmpTimeRange,
  };

  const metricPanel1: MetricPanel = {
    pluginId: 'timeseries',
    title: 'Sample Panel ONE',
    data: samplePanel,

  };
  const metricPanel2: MetricPanel = {
    pluginId: 'timeseries',
    title: 'Sample Panel TWO',
    data: samplePanel,
  };

  const metricPanel3: MetricPanel = {
    pluginId: 'timeseries',
    title: 'Sample Panel THREE',
    data: samplePanel,
  };

  console.log('samplePanel', samplePanel);

  return (
    <div>
      <button onClick={() => { resetLokiResults(); shouldRunQueries.current = true; }}>Reset Results</button>
      
      <div style={{ marginBottom: '20px' }}>
        <h3>Organized Data:</h3>
        {organizedLokiData ? (
          <pre>{JSON.stringify(organizedLokiData, null, 2)}</pre>
        ) : (
          <p>No organized data available</p>
        )}
      </div>

      <div style={{ marginBottom: '20px' }}>
        <h3>Query Data:</h3>
        {Object.keys(lokiQueryData).map((key) => (
          <React.Fragment key={key}>
            <h4>Results for process: {key}</h4>
            <pre>{JSON.stringify(lokiQueryData[key].lokiData?.series[0].fields, null, 2)}</pre>
          </React.Fragment>
        ))}
      </div>

      <div>
        <h3>Selected Rows:</h3>
        <pre>{JSON.stringify(rows, null, 2)}</pre>
      </div>
    </div>
  );
};
