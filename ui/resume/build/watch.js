"use strict";

const fs = require("fs")
const exec = require("child_process").execSync;

fs.watch("css", {}, (e, file) => {
  exec("npm run bundle-css");
});
