const Moon = require("moonjs");
const Monx = require("monx");

Moon.use(Monx)


const store = new Monx({
    state: {
        profileId: parseInt(localStorage.getItem('profileID'))
    }
})

module.exports=store;