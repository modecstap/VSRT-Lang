import React from 'react';
import ReactDOM from 'react-dom/client';
import { SWRConfig } from 'swr';
import './index.css';
import reportWebVitals from './reportWebVitals';

import { RouterProvider } from 'react-router-dom';
import { router } from './router';

const root = ReactDOM.createRoot(document.getElementById('root'));

root.render(
  <React.StrictMode>
    <SWRConfig
      value={{
        revalidateOnFocus: true,
        dedupingInterval: 2000,
      }}
    >
      <RouterProvider router={router} />
    </SWRConfig>
  </React.StrictMode>
);

reportWebVitals();
