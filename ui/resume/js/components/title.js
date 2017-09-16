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
            async create(e) {
                const store = this.get('store');
                store.dispatch(keys.BEGIN_CREATE_RESUME);
            }
        },
        store: store
    }
}

module.exports = title;
