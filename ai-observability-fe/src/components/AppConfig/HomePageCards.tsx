// // import React, { useState } from 'react';
// // import { Infrastructure, MLFrameworks, VectorDB } from 'pages';
// // import { LLM } from 'pages';

// // const HomePageCards = () => {
// //     const [selectedComponent, setSelectedComponent] = useState(null);

// //     const handleClick = (componentName: any) => {
// //         setSelectedComponent(componentName);
// //     };

// //     const renderComponent = () => {
// //         switch (selectedComponent) {
// //             case 'infrastructure':
// //                 return <Infrastructure />;
// //             case 'LLM':
// //                 return <LLM />;
// //             case 'VectorDB':
// //                 return <VectorDB />;
// //             case 'MLFrameworks':
// //                 return <MLFrameworks />;
// //             default:
// //                 return null;
// //         }
// //     };

// //     return (
// //         <div>
// //             <ul>
// //                 <li><button onClick={() => handleClick('infrastructure')}>Component 1</button></li>
// //                 <li><button onClick={() => handleClick('LLM')}>Component 2</button></li>
// //                 <li><button onClick={() => handleClick('VectorDB')}>Component 2</button></li>
// //                 <li><button onClick={() => handleClick('MLFrameworks')}>Component 2</button></li>
// //             </ul>
// //             <div>
// //                 {renderComponent()}
// //             </div>
// //         </div>
// //     );
// // };

// // export default HomePageCards;

// import React, { useState } from 'react';
// import { Infrastructure, LLM, MLFrameworks, VectorDB } from 'pages';
// import { CardElement } from '../Card/Card';

// export const HomePageCards = () => {
//   const [selectedComponent, setSelectedComponent] = useState(null);

//   const handleClick = (componentName: any) => {
//     setSelectedComponent(componentName);
//   };

//   return (
//     <div>
//       <CardElement onClick={() => handleClick('Infrastructure')} title="Component 1" description="Description for Component 1" component={<Infrastructure />} />
//       <CardElement onClick={() => handleClick('LLM')} title="Component 2" description="Description for Component 2" component={<LLM />} />
//       <CardElement onClick={() => handleClick('VectorDB')} title="Component 3" description="Description for Component 3" component={<VectorDB />} />
//       <CardElement onClick={() => handleClick('MLFrameworks')} title="Component 4" description="Description for Component 4" component={<MLFrameworks />} />
//     </div>
//   );
// };
import React, { useState } from 'react';
import { Infrastructure, LLM, VectorDB, MLFrameworks } from 'pages'; // Check this import
import { CardElement } from '../Card/Card'; // Make sure to import your CardElement component properly


// CardElement component remains the same

const HomePageCards = () => {
  const [selectedComponent, setSelectedComponent] = useState(null);

  const handleClick = (componentName: any) => {
    setSelectedComponent(componentName);
  };

  const renderComponent = () => {
    switch (selectedComponent) {
      case 'Infrastructure':
        return <Infrastructure />;
      case 'LLM':
        return <LLM />;
      case 'VectorDB':
        return <VectorDB />;
      case 'MLFrameworks':
        return <MLFrameworks />;
      default:
        return null;
    }
  };

  return (
    <div>
    <div style={{display: 'flex'}}>
      <ul style= {{display: 'flex', listStyleType: 'none'}}>
        <li>
          <CardElement
            title="Infrastructure"
            description="Monitor GPU usage"
            onClick={() => handleClick('Infrastructure')}
          />
        </li>
        <li>
          <CardElement
            title="LLM"
            description="Monitor Your LLM Usage"
            onClick={() => handleClick('LLM')}
          />
        </li>
        <li>
          <CardElement
            title="Vector DB"
            description="Monitor your DB Usage"
            onClick={() => handleClick('VectorDB')}
          />
        </li>
        <li>
          <CardElement
            title="ML Frameworks"
            description="Monitor your ML frameworks"
            onClick={() => handleClick('MLFrameworks')}
          />
        </li>
      </ul>
      </div>
      <div>
        {renderComponent()}
      </div>
    </div>
  );
};

export default HomePageCards;
