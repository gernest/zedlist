const Moon = require("moonjs");
const Monx = require("monx");
const title = require('./title.js');
const keys = require('./constants');
const bus = require('./bus');

Moon.use(Monx)

const opts = {
    state: {
        profileId: parseInt(localStorage.getItem('profileID')),
        next: 0
    },
    actions: {
        [keys.NEXT_FORM](state) {
            state.next++
            bus.emit(keys.NEXT_FORM, { next: state.next })
        }
    }
};
title.init(opts);
const store = new Monx(opts);
module.exports = store;