# OpenAI Monitoring with grafana-openai-monitoring

## Python


### Install The Dependency:  [Grafana OpenAI Monitoring on pypi](https://pypi.org/project/grafana-openai-monitoring/)

To start collecting telemetry data from your application, you first need to install the `grafana-openai-monitoring`. 

```shell
pip install grafana-openai-monitoring
```

### Create Grafana Cloud Token

A Grafana Cloud token is necessary for authentication when sending telemetry data to Grafana Cloud. The provided token should be kept secure and used in your application's environment variables for authorization.

```
```

### Paste OpenAI Token
Create an OpenAI API Key [here](https://platform.openai.com/account/api-keys) and paste the key below

```
```

### Instrument your code

To begin collecting telemetry data, initialize the `chat_v1` or `chat_v2` object at the start of your application. This simple step hooks into your application to start monitoring its performance and behavior.
    
#### Completions

```
from openai import OpenAI
from grafana_openai_monitoring import chat_v1
client = OpenAI(
    api_key="",
)
# Apply the custom decorator to the OpenAI API function
client.completions.create = chat_v1.monitor(
    client.completions.create,
    metrics_url="",  
    logs_url="",  
    metrics_username=,  
    logs_username=,  
    access_token=""
)

# Now any call to client.completions.create will be automatically tracked
response = client.completions.create(model="davinci", max_tokens=100, prompt="Isn't Grafana the best?")
print(response)
```

#### ChatCompletions

```
from openai import OpenAI
from grafana_openai_monitoring import chat_v2

client = OpenAI(
    api_key="YOUR_OPEN_AI_API_KEY",
)

# Apply the custom decorator to the OpenAI API function
client.chat.completions.create = chat_v2.monitor(
    client.chat.completions.create,
    metrics_url="",  
    logs_url="",  
    metrics_username=,  
    logs_username=,  
    access_token=""
)

# Now any call to client.chat.completions.create will be automatically tracked
response = client.chat.completions.create(model="gpt-4", max_tokens=100, messages=[{"role": "user", "content": "What is Grafana?"}])
print(response)
```

## NodeJS


### Install The Dependency:  [Grafana OpenAI Monitoring on npm](https://www.npmjs.com/package/grafana-openai-monitoring)

To start collecting telemetry data from your application, you first need to install the `grafana-openai-monitoring`. 

```shell
npm install grafana-openai-monitoring
```

### Create Grafana Cloud Token

A Grafana Cloud token is necessary for authentication when sending telemetry data to Grafana Cloud. The provided token should be kept secure and used in your application's environment variables for authorization.

```
```

### Paste OpenAI Token

Create an OpenAI API Key [here](https://platform.openai.com/account/api-keys) and paste the key below


`skabcdefghijlmopqrstuvwxyz`


### Instrument your code

To begin collecting telemetry data, initialize the `chat_v1` or `chat_v2` object at the start of your application. This simple step hooks into your application to start monitoring its performance and behavior.
    
#### Completions

```
import OpenAI from 'openai';
import { chat_v1 } from 'grafana-openai-monitoring';

const openai = new OpenAI({
apiKey: '',
});

// Patch method
chat_v1.monitor(openai, {
    metrics_url: '',
    logs_url: '',
    metrics_username: ,
    logs_username: ,
    access_token: ''
});

// Now any call to openai.completions.create will be automatically tracked
async function main() {
const completion = await openai.completions.create({
    model: 'davinci',
    max_tokens: 100,
    prompt: 'Isn't Grafana the best?',
});
console.log(completion);
}

main();
```

#### ChatCompletions
```
import OpenAI from 'openai';
import { chat_v2 } from 'grafana-openai-monitoring';

const openai = new OpenAI({
apiKey: '',
});

// Patch method
chat_v2.monitor(openai, {
metrics_url: '',
logs_url: '',
metrics_username: ,
logs_username: ,
access_token: ''
});

// Now any call to openai.chat.completions.create will be automatically tracked
async function main() {
const completion = await openai.chat.completions.create({
    model: 'gpt-4',
    max_tokens: 100,
    messages: [{ role: 'user', content: 'What is Grafana?' }],
});
console.log(completion);
}

main();
```

## Install Dashboard
Get access to pre-configured dashboard that work right away
