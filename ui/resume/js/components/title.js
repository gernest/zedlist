const http = require('../http.js');
const store = require('../store/store.js');
const keys = require('../store/constants.js');
const bus = require('../store/bus.js');

function title() {
    return {
        data() {
            return {
                title: '',
                action: 'Create',
                next: store.state.next,
                hasFlash: false,
                flashText: '',
                flashStatus: ''
            }
        },
        computed: {
            showForm: {
                get: function () {
                    return this.get('next') == 0
                }
            }
        },
        hooks: {
            init() {
                const self = this;
                bus.on(keys.NEXT_FORM, (payload) => {
                    self.set('next', payload.next)
                })
            }
        },
        methods: {
            updateTitle(e) {
                this.set('title', e.target.value)
            },
            create(e) {
                const store = this.get('store');
                store.dispatch(keys.BEGIN_CREATE_RESUME);
                const title = this.get('title');
                const self = this;
                http.sendJSON('POST', "/resume/new", {
                    profileID: store.state.profileId,
                    title: title
                })
                    .then((data) => {
                        store.dispatch(keys.CREATE_RESUME_SUCCESS, data)
                        store.dispatch(keys.NEXT_FORM)
                    }, (err) => {
                        self.set('hasFlash', true)
                        self.set('flashText', "some fish")
                        self.set('flashStatus', "error")
                        store.dispatch(keys.CREATE_RESUME_FAILED)
                    })
            }
        },
        store: store
    }
}

module.exports = title;
