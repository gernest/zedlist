const http = require('../http.js');
const store = require('../store/store.js');

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
            }
        },
        store: store
    }
}

module.exports = title;
