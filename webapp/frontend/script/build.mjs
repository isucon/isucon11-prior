#!/usr/bin/env node

import { build } from 'esbuild';
import { rm, stat, writeFile } from 'fs/promises';
// import packageJSON from './package.json.cjs';
import { dirname, join, relative } from 'path';
import { sassPlugin } from 'esbuild-sass-plugin';
import htmlPlugin from '@chialab/esbuild-plugin-html';

const __dirname = dirname(new URL(import.meta.url).pathname);
const distDir = join(__dirname, '..', 'dist');

const mode = process.env.NODE_ENV === 'production' ? 'production' : 'development';
const isProduction = mode === 'production';

const size = (path) => stat(path).then((s) => s.size);

const minify = isProduction;

/** @type {import('esbuild').BuildOptions} */
const baseOptions = {
  mainFields: ['module', 'main'],
  bundle: true,
  sourcemap: true,
  metafile: true,
  minify,
  loader: {
    '.svg': 'file',
  },
  plugins: [
    sassPlugin({
      implementation: 'node-sass',
    }),
    htmlPlugin({}),
  ],
};

const compile = async (/** @type string */ outfile, /** @type {import('esbuild').BuildOptions} */ options) => {
  const start = Date.now();
  return build({
    ...baseOptions,
    ...options,
  }).then(async (result) => {
    const complete = Date.now();
    const time = complete - start;
    await Promise.all(
      Object.keys(result.metafile.outputs).map(async (file) => {
        const fileSize = await size(file);
        console.log(`===> Complete build: ${file}`, { fileSize, minify, time });
      }),
    );
    return result;
  });
};

const browser = async () => {
  return compile('./dist/index.js', {
    platform: 'browser',
    target: ['es2018'],
    format: 'iife',
    entryPoints: ['./src/index.html'],
    outdir: distDir,
    publicPath: '/',
    entryNames: '[dir]/[name]',
    assetNames: '[dir]/[name]',
  });
};

const run = async () => {
  await rm('dist', { recursive: true, force: true });
  await browser();
  return
};

export default run();
