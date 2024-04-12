# Anthropic Monitoring with OpenLLMetry

## Python

### Install `traceloop-sdk` Library

To start collecting telemetry data from your application, you first need to install the `traceloop-sdk`. This SDK is assumed to be an OpenTelemetry-based tool that simplifies the process of collecting traces and other LLM data.

```shell
pip install traceloop-sdk
```

### Create Grafana Cloud Token

A Grafana Cloud token is necessary for authentication when sending telemetry data to Grafana Cloud. The provided token should be kept secure and used in your application's environment variables for authorization.

```
#token placeholder
```

### Add Environment Variables

Setting these environment variables configures the SDK with the necessary endpoint and headers to securely send telemetry data to Grafana Cloud.

```shell
export TRACELOOP_BASE_URL=https://otlp-gateway-prod-us-east-0.grafana.net/otlp

export TRACELOOP_HEADERS="Authorization=Basic%20ODY1MjM0OmdsY19leUp2SWpvaU5qVXlPVGt5SWl3aWJpSTZJbk4wWVdOckxUZzJOVEl6TkMxdmRHeHdMWGR5YVhSbExYZGhiR3hoYUdraUxDSnJJam9pVUZCR01EY3lTbnB2TmxKd1JqWnVaak15VmpVMFVGUTBJaXdpYlNJNmV5SnlJam9pY0hKdlpDMTFjeTFsWVhOMExUQWlmWDA9"
```

### Instrument your code

To begin collecting telemetry data, initialize the `Traceloop` object at the start of your application. This simple step hooks into your application to start monitoring its performance and behavior.
    
```python 
from traceloop.sdk import Traceloop

Traceloop.init()
```

## NodeJS

### Install Library

To start collecting telemetry data from your application, you first need to install the `traceloop-sdk`. This SDK is assumed to be an OpenTelemetry-based tool that simplifies the process of collecting traces and other LLM data.

```shell
npm install @traceloop/node-server-sdk
```

### Create Grafana Cloud Token

A Grafana Cloud token is necessary for authentication when sending telemetry data to Grafana Cloud. The provided token should be kept secure and used in your application's environment variables for authorization.

```

```

### Add Environment Variables

Setting these environment variables configures the SDK with the necessary endpoint and headers to securely send telemetry data to Grafana Cloud.

```shell
export TRACELOOP_BASE_URL=https://otlp-gateway-prod-us-east-0.grafana.net/otlp

export TRACELOOP_HEADERS="Authorization=Basic%20ODY1MjM0OmdsY19leUp2SWpvaU5qVXlPVGt5SWl3aWJpSTZJbk4wWVdOckxUZzJOVEl6TkMxdmRHeHdMWGR5YVhSbExYZGhiR3hoYUdraUxDSnJJam9pVUZCR01EY3lTbnB2TmxKd1JqWnVaak15VmpVMFVGUTBJaXdpYl
```

### Instrument your code

To begin collecting telemetry data, initialize the `Traceloop` object at the start of your application. This simple step hooks into your application to start monitoring its performance and behavior.
    
```python 
import * as traceloop from "@traceloop/node-server-sdk";

traceloop.initialize();
```

## Install Dashboard
Get access to pre-configured dashboard that work right away

