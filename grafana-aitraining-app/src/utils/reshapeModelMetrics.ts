import { DataFrame, FieldType, Field } from '@grafana/data';

interface ReshapedMetrics {
  meta: {
    startTime: string | undefined;
    endTime: string | undefined;
    sections: {
      [key: string]: string[];
    };
  };
  data: {
    [section: string]: {
      [metric: string]: DataFrame;
    };
  };
}

interface TempData {
  [key: string]: {
    [processUuid: string]: number[];
  };
}

export function reshapeModelMetrics(queryData: any): ReshapedMetrics {
  const result: ReshapedMetrics = {
    meta: { startTime: undefined, endTime: undefined, sections: {} },
    data: {},
  };

  const processUuids = Object.keys(queryData);
  const tempData: TempData = {};

  for (const processUuid of processUuids) {
    const processData = queryData[processUuid];
    const lokiData = processData.lokiData;

    if (processData.processData) {
      if (processData.processData.start_time && !result.meta.startTime) {
        result.meta.startTime = processData.processData.start_time;
      }
      if (processData.processData.end_time && !result.meta.endTime) {
        result.meta.endTime = processData.processData.end_time;
      }
    }

    if (lokiData && lokiData.series) {
      for (const series of lokiData.series) {
        if (series.fields) {
          for (const field of series.fields) {
            if (field.name === 'Line' && field.values) {
              for (let i = field.values.length - 1; i >= 0; i--) {
                const logLine = JSON.parse(field.values[i]);

                for (const key in logLine) {
                  if (logLine.hasOwnProperty(key)) {
                    const section = key.split('/')[0];
                    
                    if (!result.meta.sections[section]) {
                      result.meta.sections[section] = [];
                    }
                    if (!result.meta.sections[section].includes(key)) {
                      result.meta.sections[section].push(key);
                    }

                    if (!tempData[key]) {
                      tempData[key] = {};
                    }
                    if (!tempData[key][processUuid]) {
                      tempData[key][processUuid] = [];
                    }

                    tempData[key][processUuid].push(parseFloat(logLine[key]));
                  }
                }
              }
            }
          }
        }
      }
    }
  }

  // Convert tempData to DataFrames
  for (const section in result.meta.sections) {
    result.data[section] = {};
    for (const key of result.meta.sections[section]) {
      const fields: Field<number, number[]>[] = [
        {
          name: 'x',
          type: FieldType.number,
          values: [],
          config: {},
        }
      ];
      const maxLength = Math.max(...Object.values(tempData[key]).map(arr => arr.length));
      
      for (const processUuid in tempData[key]) {
        fields.push({
          name: processUuid,
          type: FieldType.number,
          values: tempData[key][processUuid],
          config: {},
        });
      }

      for (let i = 0; i < maxLength; i++) {
        fields[0].values.push(i);
      }

      result.data[section][key] = {
        fields,
        length: maxLength,
        refId: key,
      };
    }
  }

  return result;
}
