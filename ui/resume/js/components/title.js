
function title() {
    return {
        data() {
            return {
                title: ''
            }
        },
        methods: {
            updateTitle(e) {
                this.set('title', e.target.value)
            }
        }
    }
}

module.exports=title;
