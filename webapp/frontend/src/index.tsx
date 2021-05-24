import React from 'react';
import { render } from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { App } from './app';

const main = () => {
  const root = document.getElementById('root');
  render(
    <BrowserRouter>
      <App />
    </BrowserRouter>,
    root
  );
}

export default main();
