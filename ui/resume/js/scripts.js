/*=============================
  Primary Application Code
=============================*/

const Moon = require("moonjs");
require("./components/header.moon")(Moon);

new Moon({
  el: "#app"
});
