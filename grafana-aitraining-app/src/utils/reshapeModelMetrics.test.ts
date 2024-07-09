// reshapeModelMetrics.test.ts

import { reshapeModelMetrics } from './reshapeModelMetrics';

describe('reshapeModelMetrics', () => {
  it('should reshape and reverse model metrics correctly', () => {
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
                    '{"accuracy": 0.9338064689165435, "loss": 0.1308622321689051}',
                    '{"accuracy": 0.8730602238410783, "loss": 0.15714042533691971}',
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
                    '{"accuracy": 0.7926974470414183, "loss": 0.1273479979712903}',
                    '{"accuracy": 0.8171786292906561, "loss": 0.22870355269154058}',
                  ],
                },
              ],
            },
          ],
        },
      },
    };

    const expectedOutput = {
      meta: {
        startTime: undefined,
        endTime: undefined,
        keys: ['accuracy', 'loss'],
      },
      data: {
        accuracy: {
          '49ea35fa-e1f2-4c24-9b12-c8259b4a4446': [0.8730602238410783, 0.9338064689165435],
          '8c0254a6-363b-456b-be1f-e9ff5e8904d6': [0.8171786292906561, 0.7926974470414183],
        },
        loss: {
          '49ea35fa-e1f2-4c24-9b12-c8259b4a4446': [0.15714042533691971, 0.1308622321689051],
          '8c0254a6-363b-456b-be1f-e9ff5e8904d6': [0.22870355269154058, 0.1273479979712903],
        },
      },
    };

    const reshapedData = reshapeModelMetrics(queryData);
    expect(reshapedData).toEqual(expectedOutput);
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
        startTime: undefined,
        endTime: undefined,
        keys: [],
      },
      data: {},
    };

    const reshapedData = reshapeModelMetrics(queryData);
    expect(reshapedData).toEqual(expectedOutput);
  });
});
