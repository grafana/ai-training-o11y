import { DataFrame, FieldType, Field } from '@grafana/data';
import { config } from '@grafana/runtime';
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

// This at least looks like Grafana but unsure if the choices of colors are actually good
function generateColorPalette(count: number): string[] {
  const colors = config.theme2.visualization.palette;
  return Array.from({ length: count }, (_, i) => colors[i % colors.length]);
}

function populateTempData(queryData: any): {tempData: TempData; meta: ReshapedMetrics['meta']} {
  const tempData: TempData = {};
  const meta: ReshapedMetrics['meta'] = { startTime: undefined, endTime: undefined, sections: {} };

  for (const processUuid of Object.keys(queryData)) {
    const processData = queryData[processUuid];
    const lokiData = processData.lokiData;

    if (processData.processData) {
      if (processData.processData.start_time && !meta.startTime) {
        meta.startTime = processData.processData.start_time;
      }
      if (processData.processData.end_time && !meta.endTime) {
        meta.endTime = processData.processData.end_time;
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
                    
                    if (!meta.sections[section]) {
                      meta.sections[section] = [];
                    }
                    if (!meta.sections[section].includes(key)) {
                      meta.sections[section].push(key);
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

  return { tempData, meta };
}

export function reshapeModelMetrics(queryData: any): ReshapedMetrics {
  const processUuids = Object.keys(queryData);
  const colorPalette = generateColorPalette(processUuids.length);
  const colorMap: { [key: string]: string } = {};
  processUuids.forEach((uuid, index) => {
    colorMap[uuid] = colorPalette[index];
  });

  const { tempData, meta } = populateTempData(queryData);

  const result: ReshapedMetrics = {
    meta,
    data: {},
  };

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
          config: {
            color: {
              mode: 'fixed',
              fixedColor: colorMap[processUuid],
            },
          },
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
