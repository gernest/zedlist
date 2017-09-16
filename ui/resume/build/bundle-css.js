"use strict";

const fs = require("fs");
const path = require("path");
const Slash = require("slashjs");
const CleanCSS = require("clean-css");
const isProduction = process.env.NODE_ENV === "production";
let builtCSSPath = path.join("dist", "css", "build.min.css");

if(!fs.existsSync(builtCSSPath)) {
  fs.writeFileSync(builtCSSPath);
}

// Get all css files
const cssDir = path.join(process.cwd(), "css");
const cssFiles = fs.readdirSync(cssDir);
let css = "";

// Concat all files into one big css file
for(let i = 0; i < cssFiles.length; i++) {
  css += fs.readFileSync(path.join(cssDir, cssFiles[i]));
}

// If in production, include CSS from components
// Also Expose the Hash
if(isProduction) {
  css += fs.readFileSync(builtCSSPath);

  const hash = Slash(css).toLowerCase();
  fs.unlinkSync(builtCSSPath);
  builtCSSPath = path.join("dist", "css", `build.${hash}.css`);
  module.exports = hash;
}

// Optimize CSS
const optimizedCSS = new CleanCSS({}).minify(css);
fs.writeFileSync(builtCSSPath, optimizedCSS.styles);
