// utils/reshapeModelMetrics.ts

export function reshapeModelMetrics(queryData: any) {
  const result: any = {};

  const processUuids = Object.keys(queryData);

  console.log('[shaping] processUuids',processUuids)

  for (let i = 0; i < processUuids.length; i++) {
    const processUuid = processUuids[i];
    const processData = queryData[processUuid];
    const lokiData = processData.lokiData;

    console.log(`[shaping] process: ${processUuid}`, { processData, lokiData });

    if (lokiData && lokiData.series) {
      for (let j = 0; j < lokiData.series.length; j++) {
        const series = lokiData.series[j];

        console.log(`[shaping] series: ${j}`,series);

        if (series.fields) {
          for (let k = 0; k < series.fields.length; k++) {
            const field = series.fields[k];

            console.log(`[shaping] field: ${k}`,field);

            if (field.name === 'Line' && field.values) {
              for (let l = 0; l < field.values.length; l++) {
                const logLine = JSON.parse(field.values[l]);

                console.log(`[shaping] logLine: ${l}`,logLine);

                for (const key in logLine) {
                  if (logLine.hasOwnProperty(key)) {
                    if (!result[key]) {
                      result[key] = {};
                    }

                    if (!result[key][processUuid]) {
                      result[key][processUuid] = [];
                    }

                    result[key][processUuid].push(logLine[key]);
                  }
                }
              }
            }
          }
        }
      }
    }
  }

  console.log('[shaping] result',result);

  return result;
}
