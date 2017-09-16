const keys = require('./constants.js');

const actions = {
    [keys.BEGIN_CREATE_RESUME](state) {
        state.resume.loading = 'loading'
    },
    [keys.CREATE_RESUME_SUCCESS](state, data) {
        state.resume.loading = ''
        state.resume.action = 'Update'
        state.resume.value = data
        state.resume.created = true
    }
}

module.exports.init = (opts) => {
    opts.state.resume = {
        title: '',
        value: {},
        created: false,
        action: 'Create',
        showProgress: false,
        showForm: true,
        loading: ''
    }
    Object.assign(opts.actions, actions)
}