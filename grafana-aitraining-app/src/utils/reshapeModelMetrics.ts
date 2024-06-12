// utils/reshapeModelMetrics.ts

export function reshapeModelMetrics(queryData: any) {
  const result: any = {};

  const processUuids = Object.keys(queryData);

  for (let i = 0; i < processUuids.length; i++) {
    const processUuid = processUuids[i];
    const processData = queryData[processUuid];
    const lokiData = processData.lokiData;

    if (lokiData && lokiData.series) {
      for (let j = 0; j < lokiData.series.length; j++) {
        const series = lokiData.series[j];

        if (series.fields) {
          for (let k = 0; k < series.fields.length; k++) {
            const field = series.fields[k];

            if (field.name === 'Line' && field.values) {
              for (let l = 0; l < field.values.length; l++) {
                const logLine = JSON.parse(field.values[l]);

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

  return result;
}
