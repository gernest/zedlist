/*=============================
  Primary Application Code
=============================*/

const Moon = require("moonjs");
require("./components/title.moon")(Moon);

new Moon({
  el: "#app"
});
