const http = require('../http.js');
const store = require('../store/store.js');
const keys = require('../store/constants.js');

function title() {
    return {
        data() {
            return {
                title: '',
                action: 'Create'
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
                const self=this;
                http.sendJSON('POST', "/resume/new", {
                    profileID: store.state.profileId,
                    title: title
                })
                    .then((data) => {
                        store.dispatch(keys.CREATE_RESUME_SUCCESS, data)
                        self.build()
                    }, (err) => {
                        console.log(err);
                    })
            }
        },
        store: store
    }
}

module.exports = title;
