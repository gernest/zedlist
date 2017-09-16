const keys = require('./constants.js');

const actions = {
    [keys.BEGIN_CREATE_RESUME](state) {
        state.showForm = false
        state.showProgress = true;
    }
}

module.exports.init = (opts) => {
    opts.state.resume = {
        title: '',
        value: {},
        created: false,
        action: 'Create',
        showProgress: false,
        showForm: true
    }
    Object.assign(opts.actions, actions)
}