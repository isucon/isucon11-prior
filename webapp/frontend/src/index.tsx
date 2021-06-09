import React from 'react';
import { render } from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { App } from './app';
import { AppContextProvider } from './context';

const main = () => {
  const root = document.getElementById('root');
  render(
    <BrowserRouter>
      <AppContextProvider>
        <App />
      </AppContextProvider>
    </BrowserRouter>,
    root
  );
}

export default main();
