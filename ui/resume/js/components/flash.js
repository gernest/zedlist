

module.exports = () => {
    return {
        props: ['status', 'text'],
        data: function () {
            return {
                hasMessage: true
            }
        },
        methods: {
            destroyMessage(e) {
                this.set('hasMessage', false)
            }
        }
    }
}