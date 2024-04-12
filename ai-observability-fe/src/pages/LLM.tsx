import React, { useState } from 'react';
import Markdown from 'markdown-to-jsx';
import AnthropicMarkdown from '../instructions/llms/anthropic/openllmetry.md';
import CohereMarkdown from '../instructions/llms/cohere/openllmetry.md';
import OpenAIMarkdown from '../instructions/llms/openai/grafana-openai-monitoring.md';
import { CardElement } from 'components/Card/Card';
import InstallDashboard from 'components/InstallDashboards/InstallDashboards';
// import openaidash from '../instructions/llms/openai/dashboard.json'

export function LLM() {
  const [selectedLLM, setSelectedLLM] = useState(null);

  const handleButtonClick = (llm: any) => {
    setSelectedLLM(llm);
  };

  return (

    <div style={{ marginTop: '64px', border: '0.5px solid gray' }}>
      <h1 style={{ padding: '16px' }}>LLM Instructions</h1>
      <div style={{ display: 'flex', padding: '16px' }}>
        <CardElement
          title="OpenAI"
          description="Monitor your OpenAI Token usage"
          onClick={() => handleButtonClick('OpenAI')}
        />
        <CardElement
          title="Cohere"
          description="Monitor Your Cohere token Usage"
          onClick={() => handleButtonClick('Cohere')}
        />
        <CardElement
          title="Anthropic"
          description="Monitor Your Anthropic token Usage"
          onClick={() => handleButtonClick('Anthropic')}
        />
      </div>
      <div style={{ marginLeft: '32px', padding: '16px', marginTop: '16px' }}>

        {selectedLLM === 'Anthropic' && <Markdown disableHtml={true}>{AnthropicMarkdown}</Markdown>}
        {selectedLLM === 'Cohere' && <Markdown disableHtml={true}>{CohereMarkdown}</Markdown>}
        {selectedLLM === 'OpenAI' && <Markdown disableHtml={true}>{OpenAIMarkdown}</Markdown>}

        {selectedLLM === 'Anthropic' && <InstallDashboard filePath="https://raw.githubusercontent.com/grafana/hackathon-2024-03-tame-the-beast/main/gtm-aiobservability-app/src/instructions/llms/anthropic/dashboard.json?" />}
        {selectedLLM === 'Cohere' && <InstallDashboard filePath="https://raw.githubusercontent.com/grafana/hackathon-2024-03-tame-the-beast/main/gtm-aiobservability-app/src/instructions/llms/cohere/dashboard.json?" />}
        {selectedLLM === 'OpenAI' && <InstallDashboard filePath="https://raw.githubusercontent.com/grafana/hackathon-2024-03-tame-the-beast/main/gtm-aiobservability-app/src/instructions/llms/openai/dashboard.json?" />}

      </div>
    </div>
  );
}
