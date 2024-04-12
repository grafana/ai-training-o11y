import { chromaDirectRun, chromaDockerContainer, chromaKubernetes } from 'pages/snippets'
import React from 'react'

const Chroma = () => {
    return (
        <div>
            <h1>Monitor ChromaDB</h1>
            <h2>Understand the Environment Variables</h2>
            <p>Firstly, here’s a quick overview of what each environment variable you'll be setting does:</p>

            <ol>
                <li><code>CHROMA_OTEL_COLLECTION_ENDPOINT</code>: This is the endpoint to which the metrics from ChromaDB will be sent. It uses the OpenTelemetry protocol to communicate with Grafana Cloud.</li>
                <li><code>CHROMA_OTEL_SERVICE_NAME</code>: The name for your ChromaDB service that will appear in Grafana to identify your metrics source.</li>
                <li><code>CHROMA_OTEL_COLLECTION_HEADERS</code>: These are the headers required for authentication with your Grafana Cloud instance. This particular one includes an authorization header with basic auth credentials.
                </li>
                <li><code>CHROMA_OTEL_GRANULARITY</code>: This defines the granularity of the telemetry data collected; in this case, all data points are collected.
                </li>
            </ol>

            <h2>Configure Environment Variables</h2>
            <p>These variables need to be set in the environment where ChromaDB runs. How you do this depends on your deployment method. Here are general instructions for a few common scenarios:</p>
            <h3>Direct Run</h3>
            <p>If you're running ChromaDB directly on a host (for development or testing), you can export these variables in your terminal:</p>
            <pre>{chromaDirectRun()}</pre>
            <h3>Docker Container</h3>
            <p>For a Docker deployment, these variables can be added to your `docker run` command or your Docker Compose file:</p>
            <pre>{chromaDockerContainer()}</pre>
            <h3>Kubernetes</h3>
            <p>If you're deploying on Kubernetes, these variables should be added to your pod specifications under `env` in the deployment YAML:</p>
            <pre>{chromaKubernetes()}</pre>
            <h2>Install Dashboards</h2>
        </div>
    )
}

export default Chroma