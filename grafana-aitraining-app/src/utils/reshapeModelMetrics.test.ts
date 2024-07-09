import { reshapeModelMetrics } from './reshapeModelMetrics';

// Mock the @grafana/runtime config
jest.mock('@grafana/runtime', () => ({
  config: {
    theme2: {
      visualization: {
        palette: [
          '#7EB26D', '#EAB839', '#6ED0E0', '#EF843C', '#E24D42', '#1F78C1',
          '#BA43A9', '#705DA0', '#508642', '#CCA300', '#447EBC', '#C15C17',
          '#890F02', '#0A437C', '#6D1F62', '#584477'
        ],
      },
    },
  },
}));

describe('reshapeModelMetrics', () => {
  it('should reshape and reverse model metrics correctly with sections and DataFrames', () => {
    const queryData = {
      '49ea35fa-e1f2-4c24-9b12-c8259b4a4446': {
        processData: {},
        lokiData: {
          series: [
            {
              fields: [
                {
                  name: 'Line',
                  values: [
                    '{"train/loss": 0.5, "train/acc": 0.8, "test/loss": 0.6, "test/acc": 0.75}',
                    '{"train/loss": 0.3, "train/acc": 0.9, "test/loss": 0.4, "test/acc": 0.85}',
                  ],
                },
              ],
            },
          ],
        },
      },
      '8c0254a6-363b-456b-be1f-e9ff5e8904d6': {
        processData: {},
        lokiData: {
          series: [
            {
              fields: [
                {
                  name: 'Line',
                  values: [
                    '{"train/loss": 0.6, "train/acc": 0.75, "test/loss": 0.7, "test/acc": 0.7}',
                    '{"train/loss": 0.4, "train/acc": 0.85, "test/loss": 0.5, "test/acc": 0.8}',
                  ],
                },
              ],
            },
          ],
        },
      },
    };

    const reshapedData = reshapeModelMetrics(queryData);

    const serializedData = JSON.parse(JSON.stringify(reshapedData, (key, value) => {
      if (key === 'values' && Array.isArray(value)) {
        return Array.from(value);
      }
      return value;
    }));

    // Helper function to create expected field structure
    // eslint-disable-next-line @typescript-eslint/array-type
    const createField = (name: string, type: string, values: Array<number | string>, color?: string) => ({
      name,
      type,
      values,
      config: color ? { color: { mode: 'fixed', fixedColor: color } } : {},
    });

    // We don't know the exact colors, so we'll just check if they exist
    const color1 = serializedData.data.train['train/loss'].fields[1].config.color.fixedColor;
    const color2 = serializedData.data.train['train/loss'].fields[2].config.color.fixedColor;

    const expectedOutput = {
      meta: {
        sections: {
          train: ['train/loss', 'train/acc'],
          test: ['test/loss', 'test/acc'],
        },
      },
      data: {
        train: {
          'train/loss': {
            fields: [
              createField('x', 'number', [0, 1]),
              createField('49ea35fa-e1f2-4c24-9b12-c8259b4a4446', 'number', [0.3, 0.5], color1),
              createField('8c0254a6-363b-456b-be1f-e9ff5e8904d6', 'number', [0.4, 0.6], color2),
            ],
            length: 2,
            refId: 'train/loss',
          },
          'train/acc': {
            fields: [
              createField('x', 'number', [0, 1]),
              createField('49ea35fa-e1f2-4c24-9b12-c8259b4a4446', 'number', [0.9, 0.8], color1),
              createField('8c0254a6-363b-456b-be1f-e9ff5e8904d6', 'number', [0.85, 0.75], color2),
            ],
            length: 2,
            refId: 'train/acc',
          },
        },
        test: {
          'test/loss': {
            fields: [
              createField('x', 'number', [0, 1]),
              createField('49ea35fa-e1f2-4c24-9b12-c8259b4a4446', 'number', [0.4, 0.6], color1),
              createField('8c0254a6-363b-456b-be1f-e9ff5e8904d6', 'number', [0.5, 0.7], color2),
            ],
            length: 2,
            refId: 'test/loss',
          },
          'test/acc': {
            fields: [
              createField('x', 'number', [0, 1]),
              createField('49ea35fa-e1f2-4c24-9b12-c8259b4a4446', 'number', [0.85, 0.75], color1),
              createField('8c0254a6-363b-456b-be1f-e9ff5e8904d6', 'number', [0.8, 0.7], color2),
            ],
            length: 2,
            refId: 'test/acc',
          },
        },
      },
    };

    expect(serializedData).toEqual(expectedOutput);
  });

  it('should handle missing or empty data gracefully', () => {
    const queryData = {
      '49ea35fa-e1f2-4c24-9b12-c8259b4a4446': {
        processData: {},
        lokiData: {},
      },
      '8c0254a6-363b-456b-be1f-e9ff5e8904d6': {
        processData: {},
        lokiData: {
          series: [
            {
              fields: [],
            },
          ],
        },
      },
    };

    const expectedOutput = {
      meta: {
        sections: {},
      },
      data: {},
    };

    const reshapedData = reshapeModelMetrics(queryData);
    expect(reshapedData).toEqual(expectedOutput);
  });
});
