const Moon = require("moonjs");
const Monx = require("monx");
const title = require('./title.js')

Moon.use(Monx)

const opts = {
    state: {
        profileId: parseInt(localStorage.getItem('profileID')),
    },
    actions: {}
};
title.init(opts);
const store = new Monx(opts);
module.exports = store;