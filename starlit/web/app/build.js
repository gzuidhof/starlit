#!/usr/bin/env node

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

const { build, watch, cliopts } = require("estrella");
const { remove, copy, writeFile, mkdir } = require("fs-extra");
const chalk = require("chalk");
const sass = require("sass");

const pkg = require("./package.json");

const version = `${pkg.version}`;
const buildTime = `${new Date().toISOString()}`;

const cssEntryFile = "src/styles/index.scss";
const cssTargetFile = "dist/static/styles/app.css";

async function buildStyles(config) {
  let result = undefined;
  try {
    result = sass.renderSync({
      file: cssEntryFile,
      outputStyle: config.debug ? 'expanded' : 'compressed',
    });
  } catch (e) {
    console.error(e);
    return;
  }

  try {
    await writeFile(cssTargetFile, result.css)
    console.log(chalk.greenBright(`Wrote ${cssTargetFile}`))
  } catch (e) {
    console.warn(chalk.yellow(`Could not write ${cssTargetFile}`))
  }

  return result.stats.includedFiles;
}

(async function () {
  await remove("dist");

  await copy("static", "dist/static", { recursive: true, overwrite: true })
  await copy("static", "dist/static", { recursive: true, overwrite: true })
  await mkdir("dist/static/styles");

  // Build client bundle
  const _ = build({
    entry: "src/client.ts",
    outfile: "dist/static/js/app.js",
    bundle: true,
    define: {
      "process.env.NODE_ENV": "production",
    },
    define: {
      VERSION: `"${version}"`,
      BUILT_AT: `"${buildTime}"`,
    },
    onStart: async (config, changedFiles, context) => {
      const isInitialBuild = changedFiles.length === 0;
      if (isInitialBuild) {

        try {
          await copy("static", "dist/static", { recursive: true });
        } catch (e) {
          console.warn(chalk.yellow("Could not remove existing dist folder and copy static assets (maybe you are running wrangler dev?)"))
        }

        const cssInputFiles = await buildStyles(config);
        if (config.watch) {
          watch(cssInputFiles, f => {
            buildStyles(config);
          });
        }
      }
    },
    platform: "browser"
  });
})();
