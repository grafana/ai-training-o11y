// utils/reshapeModelMetrics.ts

interface ReshapedMetrics {
  meta: {
    startTime: string | undefined;
    endTime: string | undefined;
    keys: string[];
  };
  data: {
    [key: string]: {
      [key: string]: string[];
    };
  };
}

export function reshapeModelMetrics(queryData: any) {
  const result: ReshapedMetrics = { meta: { startTime: undefined, endTime: undefined, keys: [] }, data: {} };

  const processUuids = Object.keys(queryData);

  console.log('[shaping] processUuids', processUuids);

  for (let i = 0; i < processUuids.length; i++) {
    const processUuid = processUuids[i];
    const processData = queryData[processUuid];
    const lokiData = processData.lokiData;

    console.log(`[shaping] process: ${processUuid}`, { processData, lokiData });

    if (processData.processData) {
      if (processData.processData.start_time && !result.meta.startTime) {
        result.meta.startTime = processData.processData.start_time;
      }
      if (processData.processData.end_time && !result.meta.endTime) {
        result.meta.endTime = processData.processData.end_time;
      }
    }

    if (lokiData && lokiData.series) {
      for (let j = 0; j < lokiData.series.length; j++) {
        const series = lokiData.series[j];

        console.log(`[shaping] series: ${j}`, series);

        if (series.fields) {
          for (let k = 0; k < series.fields.length; k++) {
            const field = series.fields[k];

            console.log(`[shaping] field: ${k}`, field);

            if (field.name === 'Line' && field.values) {
              for (let l = 0; l < field.values.length; l++) {
                const logLine = JSON.parse(field.values[l]);

                console.log(`[shaping] logLine: ${l}`, logLine);

                for (const key in logLine) {
                  if (logLine.hasOwnProperty(key)) {
                    if (!result.data[key]) {
                      result.data[key] = {};
                    }

                    if (!result.data[key][processUuid]) {
                      result.data[key][processUuid] = [];
                    }

                    result.data[key][processUuid].push(logLine[key]);
                  }
                }
              }
            }
          }
        }
      }
    }
  }

  console.log('[shaping] result', result);

  result.meta.keys = Object.keys(result.data);

  return result;
}
