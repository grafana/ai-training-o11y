import React, { useState } from 'react';
import Markdown from 'markdown-to-jsx';
import ChromaMarkdown from '../instructions/vectordb/chroma.md';
import PineconeMarkdown from '../instructions/vectordb/pinecone.md';
import { CardElement } from 'components/Card/Card';
import InstallDashboard from 'components/InstallDashboards/InstallDashboards';
import Chroma from 'components/VectorDB/Chroma';
import PineCone from 'components/VectorDB/PineCone';


export function VectorDB() {
  const [selectedDB, setSelectedDB] = useState(null);

  // Function to handle button click
  const handleButtonClick = (db: any) => {
    setSelectedDB(db);
  };

  return (
      <div style={{marginTop: '64px', border: '0.5px solid gray'}}>
        <h1 style={{padding: '16px'}}>VectorDB Instructions</h1>
        <div style={{display: 'flex', padding: '16px'}}>
        <CardElement
          title="Chroma"
          description="Monitor Your LLM Usage"
          onClick={() => handleButtonClick('Chroma')}
        />
        <CardElement
          title="Pinecone"
          description="Monitor Your LLM Usage"
          onClick={() => handleButtonClick('Pinecone')}
        />
           </div>
        <div style={{ marginLeft: '32px', padding: '16px', marginTop: '16px' }}>
          {selectedDB === 'Chroma' && <Chroma />}
          {selectedDB === 'Pinecone' && <PineCone />}

          {selectedDB === 'Chroma' && <InstallDashboard filePath="https://raw.githubusercontent.com/grafana/hackathon-2024-03-tame-the-beast/main/gtm-aiobservability-app/src/instructions/vectordb/chroma.json" />}
        {selectedDB === 'Pinecone' && <InstallDashboard filePath='https://raw.githubusercontent.com/grafana/hackathon-2024-03-tame-the-beast/main/gtm-aiobservability-app/src/instructions/vectordb/pinecone.json' />}
        </div>
      </div>
  );
}

